package store_fs

import (
	"fmt"
	"os"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) CheckoutOne(
	options checkout_options.Options,
	sz *sku.Transacted,
) (col sku.CheckedOutLike, err error) {
	return s.checkoutOneNew(options, sz)
}

func (s *Store) checkoutOneNew(
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz *CheckedOut, err error) {
	cz = GetCheckedOutPool().Get()

	if err = cz.Internal.SetFromSkuLike(sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	if s.konfig.IsDryRun() {
		return
	}

	var e *KennungFDPair
	ok := false

	if e, ok = s.Get(&sz.Kennung); ok {
		var cze *External

		if cze, err = s.ReadExternalFromKennungFDPair(
			sku.ObjekteOptions{
				Mode: objekte_mode.ModeRealizeSansProto,
			},
			e,
			sz,
		); err != nil {
			if errors.Is(err, ErrExternalHasConflictMarker) && options.AllowConflicted {
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

			sku.DetermineState(cz, true)

			if !s.shouldCheckOut(options, cz, false) {
				return
			}

			if options.Path == checkout_options.PathDefault {
				if err = cz.Remove(s.standort); err != nil {
					err = errors.Wrap(err)
					return
				}
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

func (s *Store) UpdateCheckoutFromCheckedOut(
	options checkout_options.Options, // TODO CheckoutMode is currently ignored
	col sku.CheckedOutLike,
) (err error) {
	cofs := col.(*CheckedOut)
	sz := cofs.GetSku()

	var e *KennungFDPair
	ok := false

	if e, ok = s.Get(&sz.Kennung); !ok {
		return
	}

	if err = s.ReadIntoExternalFromKennungFDPair(
		sku.ObjekteOptions{
			Mode: objekte_mode.ModeRealizeSansProto,
		},
		e,
		sz,
		&cofs.External,
	); err != nil {
		if errors.Is(err, ErrExternalHasConflictMarker) && options.AllowConflicted {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	} else {
		sku.DetermineState(cofs, true)

		if !s.shouldCheckOut(options, cofs, true) {
			return
		}

		var mode checkout_mode.Mode

		if mode, err = cofs.External.GetCheckoutMode(); err != nil {
			err = errors.Wrap(err)
			return
		}

		options.CheckoutMode = mode

		if err = cofs.Remove(s.standort); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.checkoutOne(
		options,
		cofs,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) checkoutOne(
	options checkout_options.Options,
	cz *CheckedOut,
) (err error) {
	if s.konfig.IsDryRun() {
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
	inlineAkte := s.konfig.IsInlineTyp(t)

	if options.CheckoutMode.IncludesObjekte() {
		if err = cz.External.GetFDsPtr().Objekte.SetPath(filename); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if ((!inlineAkte || !options.CheckoutMode.IncludesObjekte()) &&
		!options.ForceInlineAkte) &&
		options.CheckoutMode.IncludesAkte() {

		fe := s.konfig.GetTypExtension(t.String())

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

	if err = s.fileEncoder.Encode(
		options.TextFormatterOptions,
		&cz.External,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) shouldCheckOut(
	options checkout_options.Options,
	cz *CheckedOut,
	allowMutterMatch bool,
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

	if !allowMutterMatch {
		return false
	}

	if mutter, err := s.externalStoreInfo.FuncReadSha(cz.Internal.Metadatei.Mutter()); err == nil {
		if metadatei.EqualerSansTai.Equals(&mutter.Metadatei, &cz.External.Metadatei) {
			return true
		}
	}

	return false
}

func (s *Store) filenameForTransacted(
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

func (s *Store) FileExtensionForGattung(
	gg interfaces.GattungGetter,
) string {
	return s.fileExtensions.GetFileExtensionForGattung(gg)
}
