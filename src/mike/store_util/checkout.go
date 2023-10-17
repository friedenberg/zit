package store_util

import (
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/id"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/checkout_options"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/cwd"
)

func (s *common) CheckoutQuery(
	options checkout_options.Options,
	fq matcher.FuncReaderTransactedLikePtr,
	f schnittstellen.FuncIter[*sku.CheckedOut],
) (err error) {
	if err = fq(
		func(t *sku.Transacted) (err error) {
			var cop *sku.CheckedOut

			cop, err = s.CheckoutOne(
				checkout_options.Options(options),
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
func (s *common) Checkout(
	options checkout_options.Options,
	fq matcher.FuncReaderTransactedLikePtr,
	ztw schnittstellen.FuncIter[*sku.Transacted],
) (zcs sku.CheckedOutMutableSet, err error) {
	zcs = collections_ptr.MakeMutableValueSet[sku.CheckedOut, *sku.CheckedOut](nil)
	zts := sku.MakeTransactedMutableSet()

	var l sync.Mutex

	if err = fq(
		iter.MakeChain(
			// zettel.MakeWriterKonfig(s.GetKonfig(), s.GetAkten().GetTypV0()),
			ztw,
			func(sk *sku.Transacted) (err error) {
				var z sku.Transacted

				if err = z.SetFromSkuLike(sk); err != nil {
					err = errors.Wrap(err)
					return
				}

				l.Lock()
				defer l.Unlock()

				if err = zts.AddPtr(&z); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
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

			zcs.AddPtr(zc)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s common) shouldCheckOut(
	options checkout_options.Options,
	cz *sku.CheckedOut,
) (ok bool) {
	if options.Force {
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

func (s *common) FileExtensionForGattung(
	gg schnittstellen.GattungGetter,
) string {
	return s.GetKonfig().FileExtensions.GetFileExtensionForGattung(gg)
}

func (s *common) PathForTransacted(dir string, tl *sku.Transacted) string {
	return path.Join(
		dir,
		fmt.Sprintf(
			"%s.%s",
			tl.Kennung,
			s.FileExtensionForGattung(tl),
		),
	)
}

func (s common) filenameForTransacted(
	options checkout_options.Options,
	sz *sku.Transacted,
) (originalFilename string, filename string, err error) {
	dir := s.standort.Cwd()

	switch options.Path {
	case checkout_options.PathTempLocal:
		var f *os.File

		if f, err = s.standort.FileTempLocal(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, f)

		originalFilename = f.Name()
		filename = f.Name()

		return
	default:
	}

	switch sz.GetGattung() {
	case gattung.Zettel:
		var h kennung.Hinweis

		if err = h.Set(sz.GetKennungLike().String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if originalFilename, err = id.MakeDirIfNecessary(h, dir); err != nil {
			err = errors.Wrap(err)
			return
		}

		filename = s.PathForTransacted(dir, sz)

	default:
		originalFilename = s.PathForTransacted(dir, sz)
		filename = originalFilename
	}

	return
}

func (s *common) CheckoutOne(
	options checkout_options.Options,
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

		if e, ok = s.cwdFiles.Get(sz.GetKennungLikePtr()); !ok {
			err = errors.Errorf(
				"file at %s not recognized as zettel: %s",
				filename,
				sz,
			)

			return
		}

		var cze *sku.External

		cze, err = s.ReadOneExternal(
			e,
			sz,
		)

		if err != nil {
			if errors.Is(err, sku.ErrExternalHasConflictMarker) && options.AllowConflicted {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		} else {
			cz.External = *cze
			cz.DetermineState(true)

			if !s.shouldCheckOut(options, cz) {
				return
			}
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

	if ((!inlineAkte || !options.CheckoutMode.IncludesObjekte()) && !options.ForceInlineAkte) &&
		options.CheckoutMode.IncludesAkte() {
		t := sz.GetTyp()

		fe := s.GetKonfig().TypenToExtensions[t.String()]

		if fe == "" {
			fe = t.String()
		}

		cz.External.GetFDsPtr().Akte.Path = originalFilename + "." + fe
	}

	if err = s.fileEncoder.Encode(&cz.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
