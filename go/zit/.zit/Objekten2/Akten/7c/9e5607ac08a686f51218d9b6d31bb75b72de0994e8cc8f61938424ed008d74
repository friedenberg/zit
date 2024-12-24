package store_fs

import (
	"fmt"
	"os"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) CheckoutOne(
	options checkout_options.Options,
	sz sku.TransactedGetter,
) (col sku.SkuType, err error) {
	col, _, err = s.checkoutOneIfNecessary(options, sz)
	return
}

func (s *Store) checkoutOneForReal(
	options checkout_options.Options,
	co *sku.CheckedOut,
	item *sku.FSItem,
) (err error) {
	if s.config.IsDryRun() {
		return
	}

	// delete the existing checkout if it exists in the cwd
	if options.Path == checkout_options.PathDefault {
		if err = s.RemoveItem(item); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var info checkoutFileNameInfo

	if err = s.hydrateCheckoutFileNameInfoFromCheckedOut(
		options,
		co,
		&info,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.setObjectIfNecessary(options, item, info); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.setBlobIfNecessary(
		options,
		item,
		info,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// This is necessary otherwise External is an empty sku
	sku.Resetter.ResetWith(co.GetSkuExternal(), co.GetSku())

	if err = s.WriteFSItemToExternal(item, co.GetSkuExternal()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.fileEncoder.Encode(
		options.TextFormatterOptions,
		co.GetSkuExternal(),
		item,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) setObjectIfNecessary(
	options checkout_options.Options,
	i *sku.FSItem,
	info checkoutFileNameInfo,
) (err error) {
	if !options.CheckoutMode.IncludesMetadata() {
		i.MutableSetLike.Del(&i.Object)
		i.Object.Reset()
		return
	}

	if err = i.Object.SetPath(info.objectName); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.Add(&i.Object)

	return
}

func (s *Store) setBlobIfNecessary(
	options checkout_options.Options,
	i *sku.FSItem,
	info checkoutFileNameInfo,
) (err error) {
	if info.inlineBlob && options.CheckoutMode.IncludesMetadata() ||
		options.ForceInlineBlob || !options.CheckoutMode.IncludesBlob() {
		i.MutableSetLike.Del(&i.Blob)
		i.Blob.Reset()
		return
	}

	fe := s.config.GetTypeExtension(info.tipe.String())

	if fe == "" {
		fe = info.tipe.String()
	}

	if err = i.Blob.SetPath(
		info.basename + "." + fe,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.Add(&i.Blob)

	return
}

func (s *Store) shouldCheckOut(
	options checkout_options.Options,
	cz *sku.CheckedOut,
	allowMutterMatch bool,
) bool {
	if options.Force {
		return true
	}

	eq := object_metadata.EqualerSansTai.Equals(
		&cz.GetSku().Metadata,
		&cz.GetSkuExternal().Metadata,
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
		cz.GetSku().GetObjectId().String(),
		mutter,
	); err == nil {
		if object_metadata.EqualerSansTai.Equals(&mutter.Metadata, &cz.GetSkuExternal().Metadata) {
			return true
		}
	}

	ui.Log().Print("")

	return false
}

type checkoutFileNameInfo struct {
	basename   string
	objectName string
	tipe       ids.Type
	inlineBlob bool
}

func (s *Store) hydrateCheckoutFileNameInfoFromCheckedOut(
	options checkout_options.Options,
	co *sku.CheckedOut,
	info *checkoutFileNameInfo,
) (err error) {
	if err = s.SetFilenameForTransacted(options, co.GetSku(), info); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.SetState(checked_out_state.JustCheckedOut)

	info.tipe = co.GetSku().GetType()
	info.inlineBlob = s.config.IsInlineType(info.tipe)

	return
}

func (s *Store) SetFilenameForTransacted(
	options checkout_options.Options,
	sk *sku.Transacted,
	info *checkoutFileNameInfo,
) (err error) {
	cwd := s.dirLayout.Cwd()

	if options.Path == checkout_options.PathTempLocal {
		var f *os.File

		if f, err = s.dirLayout.TempLocal.FileTemp(); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer errors.DeferredCloser(&err, f)

		info.basename = f.Name()
		info.objectName = f.Name()

		return
	}

	if sk.GetGenre() == genres.Zettel {
		var h ids.ZettelId

		if err = h.Set(sk.GetObjectId().String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if info.basename, err = id.MakeDirIfNecessary(h, cwd); err != nil {
			err = errors.Wrap(err)
			return
		}

		info.objectName = s.PathForTransacted(cwd, sk)
	} else {
		info.basename = s.PathForTransacted(cwd, sk)
		info.objectName = info.basename
	}

	return
}

func (s *Store) PathForTransacted(dir string, tl *sku.Transacted) string {
	return path.Join(
		dir,
		fmt.Sprintf(
			"%s.%s",
			&tl.ObjectId,
			s.FileExtensionForGenre(tl),
		),
	)
}

func (s *Store) FileExtensionForGenre(
	gg interfaces.GenreGetter,
) string {
	ext := s.fileExtensions.GetFileExtensionForGenre(gg)

	if ext == "" {
		panic("empty file extension")
	}

	return ext
}

func (s *Store) RemoveItem(i *sku.FSItem) (err error) {
	// TODO check conflict state
	if err = i.MutableSetLike.Each(
		func(f *fd.FD) (err error) {
			if err = f.Remove(s.dirLayout); err != nil {
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

	return
}

func (s *Store) UpdateCheckoutFromCheckedOut(
	options checkout_options.OptionsWithoutMode,
	co sku.SkuType,
) (err error) {
	o := checkout_options.Options{
		OptionsWithoutMode: options,
	}

	if o.CheckoutMode, err = s.GetCheckoutMode(
		co.GetSkuExternal(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if o.CheckoutMode == checkout_mode.None {
		return
	}

	options.Path = checkout_options.PathTempLocal

	var replacement *sku.CheckedOut
	var oldFDs, newFDs *sku.FSItem

	if oldFDs, err = s.ReadFSItemFromExternal(co.GetSkuExternal()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if replacement, newFDs, err = s.checkoutOneIfNecessary(
		o,
		co.GetSkuExternal(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer GetCheckedOutPool().Put(replacement)

	if !oldFDs.Object.IsEmpty() && !s.config.IsDryRun() {
		if err = os.Rename(
			newFDs.Object.GetPath(),
			oldFDs.Object.GetPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !oldFDs.Blob.IsEmpty() && !s.config.IsDryRun() {
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
