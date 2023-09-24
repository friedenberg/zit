package store_fs

import (
	"fmt"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/objekte_collections"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/mike/store_util"
)

func (s *Store) CheckoutQuery(
	options store_util.CheckoutOptions,
	fq objekte.FuncReaderTransactedLikePtr,
	f schnittstellen.FuncIter[*sku.CheckedOut],
) (err error) {
	if err = fq(
		func(t *sku.Transacted) (err error) {
			var cop *sku.CheckedOut

			cop, err = s.CheckoutOne(
				store_util.CheckoutOptions(options),
				t,
			)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			cop.DetermineState(true)

			if err = s.checkedOutLogPrinter(cop); err != nil {
				err = errors.Wrap(err)
				return
			}

			return f(cop)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO-P2 combine with CheckoutQuery once all matcher Query is simplified into
// just a matcher
func (s *Store) Checkout(
	options store_util.CheckoutOptions,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
	fq objekte.FuncReaderTransactedLikePtr,
	ztw schnittstellen.FuncIter[*sku.Transacted],
) (zcs schnittstellen.MutableSetLike[*sku.CheckedOut], err error) {
	zcs = collections_value.MakeMutableValueSet[*sku.CheckedOut](nil)
	zts := sku.MakeTransactedMutableSet()

	if err = fq(
		iter.MakeChain(
			zettel.MakeWriterKonfig(s.GetKonfig(), tagp),
			ztw,
			func(sk *sku.Transacted) (err error) {
				var z sku.Transacted

				if err = z.SetFromSkuLike(sk); err != nil {
					err = errors.Wrap(err)
					return
				}

				return zts.AddPtr(&z)
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = zts.EachPtr(
		func(zt *sku.Transacted) (err error) {
			var zc *sku.CheckedOut

			if zc, err = s.CheckoutOne(options, zt); err != nil {
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
	cz *sku.CheckedOut,
) (ok bool) {
	if options.Force == true {
		ok = true
		return
	}

	if cz.State == checked_out_state.StateEmpty {
		ok = true
	}

	if cz.Internal.GetMetadatei().Equals(
		cz.External.GetMetadatei(),
	) {
		return
	}

	return
}

func (s *Store) FileExtensionForGattung(
	gg schnittstellen.GattungGetter,
) string {
	return s.GetKonfig().FileExtensions.GetFileExtensionForGattung(gg)
}

func (s *Store) PathForTransacted(tl *sku.Transacted) string {
	return path.Join(
		s.Cwd(),
		fmt.Sprintf(
			"%s.%s",
			tl.Kennung,
			s.FileExtensionForGattung(tl),
		),
	)
}

func (s Store) filenameForTransacted(
	options store_util.CheckoutOptions,
	sz *sku.Transacted,
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

		filename = s.PathForTransacted(sz)

	default:
		originalFilename = s.PathForTransacted(sz)
		filename = originalFilename
	}

	return
}

func (s *Store) CheckoutOne(
	options store_util.CheckoutOptions,
	sz *sku.Transacted,
) (cz *sku.CheckedOut, err error) {
	cz = &sku.CheckedOut{}

	if err = cz.Internal.SetFromSkuLike(sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	var originalFilename, filename string

	if originalFilename, filename, err = s.filenameForTransacted(options, sz); err != nil {
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

		if cze, err = s.ReadOneExternal(
			e,
			sz,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = cz.External.SetFromSkuLike(cze); err != nil {
			err = errors.Wrap(err)
			return
		}

		cz.DetermineState(true)

		if !s.shouldCheckOut(options, cz) {
			return
		}
	}

	inlineAkte := s.GetKonfig().IsInlineTyp(sz.GetTyp())

	cz.State = checked_out_state.StateJustCheckedOut

	if err = cz.External.SetFromSkuLike(sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if options.CheckoutMode.IncludesObjekte() {
		cz.External.GetFDsPtr().Objekte.Path = filename
	}

	if (!inlineAkte || !options.CheckoutMode.IncludesObjekte()) &&
		options.CheckoutMode.IncludesAkte() {
		t := sz.GetTyp()

		fe := s.GetKonfig().TypenToExtensions[t.String()]

		if fe == "" {
			fe = t.String()
		}

		cz.External.GetFDsPtr().Akte.Path = originalFilename + "." + fe
	}

	e := objekte_collections.MakeFileEncoder(s, s.GetKonfig())

	if err = e.Encode(&cz.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
