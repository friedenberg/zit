package store_fs

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/objekte_collections"
	"github.com/friedenberg/zit/src/juliett/objekte"
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
) (zcs schnittstellen.MutableSetLike[objekte.CheckedOutLikePtr], err error) {
	zcs = collections_value.MakeMutableValueSet[objekte.CheckedOutLikePtr](nil)
	zts := collections_value.MakeMutableValueSet[sku.SkuLikePtr](nil)

	if err = s.storeObjekten.Zettel().ReadAllSchwanzen(
		iter.MakeChain(
			zettel.MakeWriterKonfig(s.erworben, s.storeObjekten.Typ()),
			ztw,
			func(sk sku.SkuLikePtr) (err error) {
				var z sku.Transacted

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
		func(zt sku.SkuLikePtr) (err error) {
			var zc objekte.CheckedOutLikePtr

			if zc, err = s.CheckoutOneZettel(options, zt); err != nil {
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
	cz objekte.CheckedOutLikePtr,
) (ok bool) {
	if options.Force == true {
		ok = true
		return
	}

	if cz.GetState() == checked_out_state.StateEmpty {
		ok = true
	}

	if cz.GetState() == checked_out_state.StateEmpty {
		ok = true
	}

	if cz.GetInternalLike().GetMetadatei().Equals(
		cz.GetExternalLike().GetMetadatei(),
	) {
		return
	}

	return
}

func (s Store) filenameForZettelTransacted(
	options store_util.CheckoutOptions,
	sz sku.SkuLikePtr,
) (originalFilename string, filename string, err error) {
	switch sz.GetGattung() {
	case gattung.Zettel:
		var h kennung.Hinweis

		if err = h.Set(sz.GetKennungLike().String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if originalFilename, err = id.MakeDirIfNecessary(h, s.Cwd()); err != nil {
			err = errors.Wrap(err)
			return
		}

		filename = originalFilename + s.erworben.GetZettelFileExtension()

	default:
		originalFilename = sz.GetKennungLike().String() + "." + s.erworben.FileExtensions.GetFileExtensionForGattung(
			sz.GetKennungLike(),
		)

		filename = originalFilename
	}

	return
}

func (s *Store) checkoutOneGeneric(
	options store_util.CheckoutOptions,
	t sku.SkuLikePtr,
) (cop objekte.CheckedOutLikePtr, err error) {
	switch tt := t.(type) {
	case *sku.Transacted:
		cop, err = s.CheckoutOneZettel(options, tt)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		cop, err = s.storeObjekten.CheckoutOne(store_util.CheckoutOptions(options), tt)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	cop.DetermineState(true)

	if err = s.checkedOutLogPrinter(cop); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CheckoutOneZettel(
	options store_util.CheckoutOptions,
	sz sku.SkuLikePtr,
) (cz objekte.CheckedOutLikePtr, err error) {
	cz = &objekte.CheckedOut2{}

	if err = cz.GetInternalLikePtr().SetFromSkuLike(sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	var originalFilename, filename string

	if originalFilename, filename, err = s.filenameForZettelTransacted(options, sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if files.Exists(filename) {
		var e *cwd.Zettel
		ok := false

		if e, ok = options.Cwd.Get(sz.GetKennungLikePtr()); !ok {
			err = errors.Errorf(
				"file at %s not recognized as zettel: %s",
				filename,
				sz,
			)
			return
		}

		var cze objekte.ExternalLikePtr

		if cze, err = s.storeObjekten.ReadOneExternal(
			e,
			sz,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = cz.GetExternalLikePtr().SetFromSkuLike(cze); err != nil {
			err = errors.Wrap(err)
			return
		}

		cz.DetermineState(true)

		if !s.shouldCheckOut(options, cz) {
			return
		}
	}

	inlineAkte := s.erworben.IsInlineTyp(sz.GetTyp())

	cz.SetState(checked_out_state.StateJustCheckedOut)

	if err = cz.GetExternalLikePtr().SetFromSkuLike(sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if options.CheckoutMode.IncludesObjekte() {
		cz.GetExternalLikePtr().GetFDsPtr().Objekte.Path = filename
	}

	if (!inlineAkte || !options.CheckoutMode.IncludesObjekte()) &&
		options.CheckoutMode.IncludesAkte() {
		t := sz.GetTyp()

		fe := s.erworben.TypenToExtensions[t.String()]

		if fe == "" {
			fe = t.String()
		}

		cz.GetExternalLikePtr().GetFDsPtr().Akte.Path = originalFilename + "." + fe
	}

	e := objekte_collections.MakeFileEncoder(
		s.storeObjekten,
		s.erworben,
	)

	if err = e.Encode(cz.GetExternalLikePtr()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
