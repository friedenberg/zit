package store_fs

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kilo/cwd"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
	"github.com/friedenberg/zit/src/lima/store_objekten"
)

func (s *Store) CheckoutQuery(
	options CheckoutOptions,
	ms kennung.MetaSet,
	f schnittstellen.FuncIter[objekte.CheckedOutLike],
) (err error) {
	if err = s.storeObjekten.Query(
		ms,
		func(t objekte.TransactedLikePtr) (err error) {
			var co objekte.CheckedOutLikePtr

			if co, err = s.checkoutOneGeneric(options, t); err != nil {
				err = errors.Wrap(err)
				return
			}

			return f(co)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Checkout(
	options CheckoutOptions,
	ztw schnittstellen.FuncIter[*zettel.Transacted],
) (zcs zettel.MutableSetCheckedOut, err error) {
	zcs = zettel.MakeMutableSetCheckedOutUnique(0)
	zts := zettel.MakeMutableSetUnique(0)

	if err = s.storeObjekten.Zettel().ReadAllSchwanzen(
		iter.MakeChain(
			zettel.MakeWriterKonfig(s.erworben),
			ztw,
			collections.AddClone[zettel.Transacted, *zettel.Transacted](zts),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = zts.Each(
		func(zt *zettel.Transacted) (err error) {
			var zc zettel.CheckedOut

			if zc, err = s.CheckoutOneZettel(options, *zt); err != nil {
				err = errors.Wrap(err)
				return
			}

			zcs.Add(zc)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Store) shouldCheckOut(
	options CheckoutOptions,
	cz zettel.CheckedOut,
) (ok bool) {
	switch {
	case cz.Internal.GetMetadatei().Equals(cz.External.GetMetadatei()):
		cz.State = objekte.CheckedOutStateJustCheckedOutButSame

	case options.Force || cz.State == objekte.CheckedOutStateEmpty:
		ok = true
	}

	return
}

func (s Store) filenameForZettelTransacted(
	options CheckoutOptions,
	sz zettel.Transacted,
) (originalFilename string, filename string, err error) {
	if originalFilename, err = id.MakeDirIfNecessary(sz.Sku.GetKennung(), s.Cwd()); err != nil {
		err = errors.Wrap(err)
		return
	}

	filename = originalFilename + s.erworben.GetZettelFileExtension()

	return
}

func (s *Store) checkoutOneGeneric(
	options CheckoutOptions,
	t objekte.TransactedLike,
) (cop objekte.CheckedOutLikePtr, err error) {
	switch tt := t.(type) {
	case *zettel.Transacted:
		var co zettel.CheckedOut
		co, err = s.CheckoutOneZettel(options, *tt)
		cop = &co

	case *kasten.Transacted:
		cop, err = s.storeObjekten.Kasten().CheckoutOne(store_objekten.CheckoutOptions(options), tt)

	case *typ.Transacted:
		cop, err = s.storeObjekten.Typ().CheckoutOne(store_objekten.CheckoutOptions(options), tt)

	case *etikett.Transacted:
		cop, err = s.storeObjekten.Etikett().CheckoutOne(store_objekten.CheckoutOptions(options), tt)

	default:
		// err = errors.Implement()
		err = gattung.MakeErrUnsupportedGattung(tt.GetSku())
		return
	}

	cop.DetermineState()
	s.checkedOutLogPrinter(cop)

	return
}

func (s *Store) CheckoutOneZettel(
	options CheckoutOptions,
	sz zettel.Transacted,
) (cz zettel.CheckedOut, err error) {
	cz.Internal = sz

	var originalFilename, filename string

	if originalFilename, filename, err = s.filenameForZettelTransacted(options, sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if files.Exists(filename) {
		var e cwd.Zettel
		ok := false

		if e, ok = options.Cwd.GetZettel(sz.Sku.GetKennung()); !ok {
			err = errors.Errorf(
				"file at %s not recognized as zettel: %s",
				filename,
				sz,
			)
			return
		}

		if cz.External, err = s.storeObjekten.Zettel().ReadOneExternal(
			e,
			&sz,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		cz.DetermineState()

		if !s.shouldCheckOut(options, cz) {
			return
		}
	}

	inlineAkte := s.erworben.IsInlineTyp(sz.GetTyp())

	cz.State = objekte.CheckedOutStateJustCheckedOut
	cz.External = zettel.External{
		Akte: sz.Akte,
		Sku: sku.External[kennung.Hinweis, *kennung.Hinweis]{
			WithKennung: sku.WithKennung[kennung.Hinweis, *kennung.Hinweis]{
				Kennung:   sz.Sku.GetKennung(),
				Metadatei: sz.GetMetadatei(),
			},
			ObjekteSha: sz.Sku.ObjekteSha,
		},
	}

	if options.CheckoutMode.IncludesObjekte() {
		cz.External.Sku.FDs.Objekte.Path = filename
	}

	if (!inlineAkte || !options.CheckoutMode.IncludesObjekte()) &&
		options.CheckoutMode.IncludesAkte() {
		t := sz.GetTyp()

		ty := s.erworben.GetApproximatedTyp(t).ApproximatedOrActual()

		var fe string

		if ty != nil {
			fe = ty.Akte.FileExtension
		}

		if fe == "" {
			fe = t.String()
		}

		cz.External.Sku.FDs.Akte.Path = originalFilename + "." + fe
	}

	e := zettel_external.MakeFileEncoder(
		s.storeObjekten,
		s.erworben,
	)

	if err = e.Encode(&cz.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
