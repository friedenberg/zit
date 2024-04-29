package store

import (
	"fmt"
	"os"
	"path"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/src/bravo/id"
	"code.linenisgreat.com/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/kilo/cwd"
)

func (s *Store) CheckoutQuery(
	options checkout_options.Options,
	qg *query.Group,
	f schnittstellen.FuncIter[*sku.CheckedOut],
) (err error) {
	if err = s.QueryWithCwd(
		qg,
		func(t *sku.Transacted) (err error) {
			var cop *sku.CheckedOut

			if cop, err = s.CheckoutOne(
				checkout_options.Options(options),
				t,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			cop.DetermineState(true)

			if err = s.checkedOutLogPrinter(cop); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = f(cop); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Store) shouldCheckOut(
	options checkout_options.Options,
	cz *sku.CheckedOut,
) bool {
	if options.Force {
		return true
	}

	if cz.State == checked_out_state.StateEmpty {
		return true
	}

	eq := metadatei.EqualerSansTai.Equals(&cz.Internal.Metadatei, &cz.External.Metadatei)

	if eq {
		return true
	}

	return false
}

func (s *Store) FileExtensionForGattung(
	gg schnittstellen.GattungGetter,
) string {
	return s.GetKonfig().FileExtensions.GetFileExtensionForGattung(gg)
}

func (s *Store) PathForTransacted(dir string, tl *sku.Transacted) string {
	return path.Join(
		dir,
		fmt.Sprintf(
			"%s.%s",
			&tl.Kennung,
			s.FileExtensionForGattung(tl),
		),
	)
}

func (s Store) filenameForTransacted(
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

		if err = h.Set(sz.GetKennung().String()); err != nil {
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

func (s *Store) UpdateCheckoutOne(
	options checkout_options.Options, // TODO CheckoutMode is currently ignored
	sz *sku.Transacted,
) (cz *sku.CheckedOut, err error) {
	var e *cwd.Zettel
	ok := false

	if e, ok = s.cwdFiles.Get(&sz.Kennung); !ok {
		return
	}

	cz = &sku.CheckedOut{}

	if err = cz.Internal.SetFromSkuLike(sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	var cze *sku.External

	if cze, err = s.ReadOneExternal(
		e,
		sz,
	); err != nil {
		if errors.Is(err, sku.ErrExternalHasConflictMarker) && options.AllowConflicted {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = cz.External.SetFromSkuLike(cze); err != nil {
			err = errors.Wrap(err)
			return
		}

		cz.DetermineState(true)

		if !s.shouldCheckOut(options, cz) {
			return
		}

		if err = cz.Remove(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var mode checkout_mode.Mode

	if mode, err = cz.External.GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	options.CheckoutMode = mode

	if err = s.checkoutOne(
		options,
		cz,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CheckoutOne(
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz *sku.CheckedOut, err error) {
	cz = &sku.CheckedOut{}

	if err = cz.Internal.SetFromSkuLike(sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.GetKonfig().DryRun {
		return
	}

	var e *cwd.Zettel
	ok := false

	if e, ok = s.cwdFiles.Get(&sz.Kennung); ok {
		var cze *sku.External

		if cze, err = s.ReadOneExternal(
			e,
			sz,
		); err != nil {
			if errors.Is(err, sku.ErrExternalHasConflictMarker) && options.AllowConflicted {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		} else {
			if err = cz.External.SetFromSkuLike(cze); err != nil {
				err = errors.Wrap(err)
				return
			}

			cz.DetermineState(true)

			if !s.shouldCheckOut(options, cz) {
				return
			}

			if err = cz.Remove(); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if err = s.checkoutOne(
		options,
		cz,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) checkoutOne(
	options checkout_options.Options,
	cz *sku.CheckedOut,
) (err error) {
	if s.GetKonfig().DryRun {
		return
	}

	var originalFilename, filename string

	if originalFilename, filename, err = s.filenameForTransacted(
		options,
		&cz.Internal,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	cz.State = checked_out_state.StateJustCheckedOut

	t := cz.Internal.GetTyp()
	inlineAkte := s.GetKonfig().IsInlineTyp(t)

	if options.CheckoutMode.IncludesObjekte() {
		if err = cz.External.GetFDsPtr().Objekte.SetPath(filename); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if ((!inlineAkte || !options.CheckoutMode.IncludesObjekte()) &&
		!options.ForceInlineAkte) &&
		options.CheckoutMode.IncludesAkte() {

		fe := s.GetKonfig().TypenToExtensions[t.String()]

		if fe == "" {
			fe = t.String()
		}

		if err = cz.External.GetFDsPtr().Akte.SetPath(
			originalFilename + "." + fe,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = cz.External.SetFromSkuLike(&cz.Internal); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.fileEncoder.Encode(&cz.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
