package store_fs

import (
	"fmt"
	"os"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) CheckoutOne(
	options checkout_options.Options,
	sz *sku.Transacted,
) (col sku.CheckedOutLike, err error) {
	col, _, err = s.checkoutOneNew(options, sz)
	return
}

func (s *Store) checkoutOneNew(
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz *CheckedOut, fds *Item, err error) {
	cz = GetCheckedOutPool().Get()
	cz.External.item.Reset()

	sku.Resetter.ResetWith(&cz.Internal, sz)

	if s.config.IsDryRun() {
		return
	}

	var i *Item
	ok := false

	ui.TodoP4("cleanup")
	if i, ok = s.Get(&sz.ObjectId); ok {
		ui.Log().Print("E", i.MutableSetLike)
		var cze *External

		if cze, err = s.ReadExternalFromItem(
			sku.CommitOptions{
				Mode: objekte_mode.ModeRealizeSansProto,
			},
			i,
			sz,
		); err != nil {
			if errors.Is(err, ErrExternalHasConflictMarker) && options.AllowConflicted {
				cz.External.ResetWith(cze)
				sku.DetermineState(cz, true)
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		} else {
			cz.External.ResetWith(cze)

			sku.DetermineState(cz, true)

			if !s.shouldCheckOut(options, cz, true) {
				ui.Log().Print("")
				return
			}
			ui.Log().Print("")

			if options.Path == checkout_options.PathDefault {
				if err = s.Remove(cz); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
			ui.Log().Print("")
		}
	}

	ui.Log().Print("")

	if fds, err = s.checkoutOne(
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

	if options.CheckoutMode, err = s.GetCheckoutModeOrError(
		col.GetSkuExternalLike(),
		options.CheckoutMode,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	options.Path = checkout_options.PathTempLocal

	var replacement *CheckedOut
	var oldFDs, newFDs *Item

	if oldFDs, err = s.ReadFromExternal(col.GetSkuExternalLike()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if replacement, newFDs, err = s.checkoutOneNew(
		options,
		cofs.GetSkuExternalLike().GetSku(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer GetCheckedOutPool().Put(replacement)

	if !oldFDs.Object.IsEmpty() {
		if err = os.Rename(
			newFDs.Object.GetPath(),
			oldFDs.Object.GetPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !oldFDs.Blob.IsEmpty() {
		if err = os.Rename(
			newFDs.Blob.GetPath(),
			oldFDs.Blob.GetPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) checkoutOne(
	options checkout_options.Options,
	cz *CheckedOut,
) (i *Item, err error) {
	if i, err = s.ReadFromExternal(&cz.External); err != nil {
		err = errors.Wrap(err)
		return
	}

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

	cz.State = checked_out_state.JustCheckedOut

	t := cz.Internal.GetType()
	inlineBlob := s.config.IsInlineType(t)

	if options.CheckoutMode.IncludesMetadata() {
		if err = i.Object.SetPath(filename); err != nil {
			err = errors.Wrap(err)
			return
		}

		i.Add(&i.Object)
	} else {
		i.MutableSetLike.Del(&i.Object)
		i.Object.Reset()
	}

	if ((!inlineBlob || !options.CheckoutMode.IncludesMetadata()) &&
		!options.ForceInlineBlob) &&
		options.CheckoutMode.IncludesBlob() {

		fe := s.config.GetTypeExtension(t.String())

		if fe == "" {
			fe = t.String()
		}

		if err = i.Blob.SetPath(
			originalFilename + "." + fe,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		i.Add(&i.Blob)
	} else {
		i.MutableSetLike.Del(&i.Blob)
		i.Blob.Reset()
	}

	sku.Resetter.ResetWith(&cz.External, &cz.Internal)

	if err = s.fileEncoder.Encode(
		options.TextFormatterOptions,
		&cz.External,
		i,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.WriteToExternal(i, &cz.External); err != nil {
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

	eq := object_metadata.EqualerSansTai.Equals(
		&cz.Internal.Metadata,
		&cz.External.Transacted.Metadata,
	)

	if eq {
		return true
	}

	if !allowMutterMatch {
		ui.Log().Print("")
		return false
	}

	mutter := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(mutter)

	if err := s.externalStoreSupplies.FuncReadOneInto(
		cz.Internal.GetObjectId().String(),
		mutter,
	); err == nil {
		if object_metadata.EqualerSansTai.Equals(&mutter.Metadata, &cz.External.Transacted.Metadata) {
			return true
		}
	}

	ui.Log().Print("")

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

func (s *Store) Remove(e *CheckedOut) (err error) {
	var i *Item

	if i, err = s.ReadFromExternal(&e.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO check conflict state
	if err = i.MutableSetLike.Each(
		func(f *fd.FD) (err error) {
			if err = f.Remove(s.fs_home); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.Reset()

	if err = s.WriteToExternal(i, &e.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
