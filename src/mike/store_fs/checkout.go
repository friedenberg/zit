package store_fs

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/objekte_collections"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/checked_out"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/mike/store_util"
)

func (s *Store) CheckoutQuery(
	options store_util.CheckoutOptions,
	ms matcher.Query,
	f schnittstellen.FuncIter[objekte.CheckedOutLike],
) (err error) {
	if err = s.storeObjekten.Query(
		ms,
		func(t sku.SkuLikePtr) (err error) {
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
	options store_util.CheckoutOptions,
	ztw schnittstellen.FuncIter[sku.SkuLikePtr],
) (zcs zettel.MutableSetCheckedOut, err error) {
	zcs = zettel.MakeMutableSetCheckedOutUnique(0)
	zts := zettel.MakeMutableSetUnique(0)

	if err = s.storeObjekten.Zettel().ReadAllSchwanzen(
		iter.MakeChain(
			zettel.MakeWriterKonfig(s.erworben, s.storeObjekten.Typ()),
			ztw,
			func(sk sku.SkuLikePtr) (err error) {
				var z transacted.Zettel

				if err = z.SetFromSkuLike(sk); err != nil {
					err = errors.Wrap(err)
					return
				}

				return zts.Add(&z)
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = zts.Each(
		func(zt *transacted.Zettel) (err error) {
			var zc checked_out.Zettel

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
	options store_util.CheckoutOptions,
	cz checked_out.Zettel,
) (ok bool) {
	switch {
	case cz.Internal.GetMetadatei().Equals(cz.External.GetMetadatei()):
		cz.State = checked_out_state.StateJustCheckedOut

	case options.Force || cz.State == checked_out_state.StateEmpty:
		ok = true
	}

	return
}

func (s Store) filenameForZettelTransacted(
	options store_util.CheckoutOptions,
	sz transacted.Zettel,
) (originalFilename string, filename string, err error) {
	if originalFilename, err = id.MakeDirIfNecessary(sz.GetKennung(), s.Cwd()); err != nil {
		err = errors.Wrap(err)
		return
	}

	filename = originalFilename + s.erworben.GetZettelFileExtension()

	return
}

func (s *Store) checkoutOneGeneric(
	options store_util.CheckoutOptions,
	t sku.SkuLike,
) (cop objekte.CheckedOutLikePtr, err error) {
	switch tt := t.(type) {
	case *transacted.Zettel:
		var co checked_out.Zettel
		co, err = s.CheckoutOneZettel(options, *tt)
		cop = &co

	case *transacted.Kasten:
		cop, err = s.storeObjekten.Kasten().CheckoutOne(store_util.CheckoutOptions(options), tt)

	case *transacted.Typ:
		cop, err = s.storeObjekten.CheckoutOne(store_util.CheckoutOptions(options), tt)

	case *sku.Transacted2:
		cop, err = s.storeObjekten.CheckoutOne(store_util.CheckoutOptions(options), tt)

	case *transacted.Etikett:
		cop, err = s.storeObjekten.Etikett().CheckoutOne(store_util.CheckoutOptions(options), tt)

	default:
		// err = errors.Implement()
		err = gattung.MakeErrUnsupportedGattung(tt.GetSkuLike())
		return
	}

	cop.DetermineState(true)
	s.checkedOutLogPrinter(cop)

	return
}

func (s *Store) CheckoutOneZettel(
	options store_util.CheckoutOptions,
	sz transacted.Zettel,
) (cz checked_out.Zettel, err error) {
	cz.Internal = sz

	var originalFilename, filename string

	if originalFilename, filename, err = s.filenameForZettelTransacted(options, sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if files.Exists(filename) {
		var e *cwd.Zettel
		ok := false

		if e, ok = options.Cwd.GetZettel(sz.GetKennungPtr()); !ok {
			err = errors.Errorf(
				"file at %s not recognized as zettel: %s",
				filename,
				sz,
			)
			return
		}

		if cz.External, err = s.storeObjekten.Zettel().ReadOneExternal(
			*e,
			&sz,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		cz.DetermineState(true)

		if !s.shouldCheckOut(options, cz) {
			return
		}
	}

	inlineAkte := s.erworben.IsInlineTyp(sz.GetTyp())

	cz.State = checked_out_state.StateJustCheckedOut
	cz.External = sz.GetExternal()

	if options.CheckoutMode.IncludesObjekte() {
		cz.External.FDs.Objekte.Path = filename
	}

	if (!inlineAkte || !options.CheckoutMode.IncludesObjekte()) &&
		options.CheckoutMode.IncludesAkte() {
		t := sz.GetTyp()

		fe := s.erworben.TypenToExtensions[t.String()]

		if fe == "" {
			fe = t.String()
		}

		cz.External.FDs.Akte.Path = originalFilename + "." + fe
	}

	e := objekte_collections.MakeFileEncoder(
		s.storeObjekten,
		s.erworben,
	)

	if err = e.Encode(&cz.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
