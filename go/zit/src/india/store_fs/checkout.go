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
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
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

	if s.config.IsDryRun() {
		return
	}

	var e *ObjectIdFDPair
	ok := false

	if e, ok = s.Get(&sz.ObjectId); ok {
		var cze *External

		if cze, err = s.ReadExternalFromObjectIdFDPair(
			sku.CommitOptions{
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
				if err = cz.Remove(s.fs_home); err != nil {
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

	var e *ObjectIdFDPair
	ok := false

	if e, ok = s.Get(&sz.ObjectId); !ok {
		return
	}

	if err = s.ReadIntoExternalFromObjectIdFDPair(
		sku.CommitOptions{
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

		if err = cofs.Remove(s.fs_home); err != nil {
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
	if s.config.IsDryRun() {
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

	t := cz.Internal.GetType()
	inlineBlob := s.config.IsInlineType(t)

	if options.CheckoutMode.IncludesMetadata() {
		if err = cz.External.GetFDsPtr().Object.SetPath(filename); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if ((!inlineBlob || !options.CheckoutMode.IncludesMetadata()) &&
		!options.ForceInlineBlob) &&
		options.CheckoutMode.IncludesBlob() {

		fe := s.config.GetTypeExtension(t.String())

		if fe == "" {
			fe = t.String()
		}

		if err = cz.External.GetFDsPtr().Blob.SetPath(
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

	eq := object_metadata.EqualerSansTai.Equals(&cz.Internal.Metadata, &cz.External.Metadata)

	if eq {
		return true
	}

	if !allowMutterMatch {
		return false
	}

	if mutter, err := s.externalStoreInfo.FuncReadSha(cz.Internal.Metadata.Mutter()); err == nil {
		if object_metadata.EqualerSansTai.Equals(&mutter.Metadata, &cz.External.Metadata) {
			return true
		}
	}

	return false
}

func (s *Store) filenameForTransacted(
	options checkout_options.Options,
	sz *sku.Transacted,
) (originalFilename string, filename string, err error) {
	dir := s.fs_home.Cwd()

	switch options.Path {
	case checkout_options.PathTempLocal:
		var f *os.File

		if f, err = s.fs_home.FileTempLocal(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, f)

		originalFilename = f.Name()
		filename = f.Name()

		return
	default:
	}

	switch sz.GetGenre() {
	case genres.Zettel:
		var h ids.ZettelId

		if err = h.Set(sz.GetObjectId().String()); err != nil {
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
			&tl.ObjectId,
			s.FileExtensionForGattung(tl),
		),
	)
}

func (s *Store) FileExtensionForGattung(
	gg interfaces.GenreGetter,
) string {
	return s.fileExtensions.GetFileExtensionForGattung(gg)
}
