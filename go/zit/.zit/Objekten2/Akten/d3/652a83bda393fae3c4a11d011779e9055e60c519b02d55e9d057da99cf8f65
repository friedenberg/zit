package store_fs

import (
	"bufio"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) Merge(conflicted sku.Conflicted) (err error) {
	var original *sku.FSItem

	if original, err = s.ReadFSItemFromExternal(
		conflicted.CheckedOut.GetSkuExternal(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var skuReplacement *sku.Transacted

	if skuReplacement, err = s.MakeMergedTransacted(conflicted); err != nil {
		if sku.IsErrMergeConflict(err) {
			if err = s.GenerateConflictMarker(
				conflicted,
				conflicted.CheckedOut,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	if original.Len() == 0 {
		// generate check out item
		// TODO if original is empty, it means this was not a checked out conflict but
		// a remote conflict
	}

	var replacement *sku.FSItem

	if replacement, err = s.ReadFSItemFromExternal(skuReplacement); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !original.Object.IsEmpty() && !replacement.Object.IsEmpty() {
		if err = files.Rename(
			replacement.Object.GetPath(),
			original.Object.GetPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !original.Blob.IsEmpty() && !replacement.Blob.IsEmpty() {
		if err = files.Rename(
			replacement.Blob.GetPath(),
			original.Blob.GetPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) checkoutConflictedForMerge(
	tm sku.Conflicted,
	mode checkout_mode.Mode,
) (local, base, remote *sku.FSItem, err error) {
	if _, local, err = s.checkoutOneForMerge(mode, tm.Local); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, base, err = s.checkoutOneForMerge(mode, tm.Base); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, remote, err = s.checkoutOneForMerge(mode, tm.Remote); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) MakeMergedTransacted(
	conflicted sku.Conflicted,
) (merged *sku.Transacted, err error) {
	if err = conflicted.MergeTags(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var localItem, baseItem, remoteItem *sku.FSItem

	inlineBlob := conflicted.IsAllInlineType(s.config)

	mode := checkout_mode.MetadataAndBlob

	if !inlineBlob {
		mode = checkout_mode.MetadataOnly
	}

	if localItem, baseItem, remoteItem, err = s.checkoutConflictedForMerge(
		conflicted,
		mode,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var mergedItem *sku.FSItem
	var diff3Error error

	mergedItem, diff3Error = s.runDiff3(
		localItem,
		baseItem,
		remoteItem,
	)

	if diff3Error != nil {
		err = errors.Wrap(diff3Error)
		return
	}

	// ui.Debug().Print(
	// "merged", mergedItem.Debug(),
	// "local", localItem.Debug(),
	// "base", baseItem.Debug(),
	// "remote", remoteItem.Debug(),
	// )

	localItem.ResetWith(mergedItem)

	merged = GetExternalPool().Get()

	merged.ObjectId.ResetWith(&conflicted.GetSku().ObjectId)

	if err = s.WriteFSItemToExternal(localItem, merged); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.HydrateExternalFromItem(
		sku.CommitOptions{
			Mode: object_mode.ModeUpdateTai,
		},
		mergedItem,
		conflicted.GetSku(),
		conflicted.GetSkuExternal(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) checkoutOneForMerge(
	mode checkout_mode.Mode,
	sz *sku.Transacted,
) (cz *sku.CheckedOut, i *sku.FSItem, err error) {
	if sz == nil {
		i = &sku.FSItem{}
		i.Reset()
		return
	}

	options := checkout_options.Options{
		CheckoutMode: mode,
		OptionsWithoutMode: checkout_options.OptionsWithoutMode{
			AllowConflicted: true,
			Path:            checkout_options.PathTempLocal,
			ForceInlineBlob: true,
			Force:           true,
		},
	}

	cz = GetCheckedOutPool().Get()
	sku.Resetter.ResetWith(cz.GetSku(), sz)

	if i, err = s.ReadFSItemFromExternal(cz.GetSkuExternal()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.checkoutOneForReal(
		options,
		cz,
		i,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.WriteFSItemToExternal(i, cz.GetSkuExternal()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) GenerateConflictMarker(
	conflicted sku.Conflicted,
	co *sku.CheckedOut,
) (err error) {
	var f *os.File

	if f, err = s.dirLayout.TempLocal.FileTemp(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	bw := bufio.NewWriter(f)
	defer errors.DeferredFlusher(&err, bw)

	blobStore := s.externalStoreSupplies.BlobStore.GetInventoryList()

	if _, err = blobStore.WriteBlobToWriter(
		builtin_types.DefaultOrPanic(genres.InventoryList),
		conflicted,
		bw,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var i *sku.FSItem

	if i, err = s.ReadFSItemFromExternal(co.GetSkuExternal()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.GenerateConflictFD(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if co.GetSkuExternal().GetGenre() == genres.Zettel {
		var h ids.ZettelId

		if err = h.Set(co.GetSkuExternal().GetObjectId().String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = id.MakeDirIfNecessary(
			h,
			s.dirLayout.Cwd(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = os.Rename(
		f.Name(),
		i.Conflict.GetPath(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.SetState(checked_out_state.Conflicted)

	return
}

func (s *Store) RunMergeTool(
	tool []string,
	conflicted sku.Conflicted,
) (co *sku.CheckedOut, err error) {
	if len(tool) == 0 {
		err = errors.Errorf("no utility provided")
		return
	}

	co = conflicted.CheckedOut

	inlineBlob := conflicted.IsAllInlineType(s.config)

	mode := checkout_mode.MetadataAndBlob

	if !inlineBlob {
		mode = checkout_mode.MetadataOnly
	}

	var leftItem, middleItem, rightItem *sku.FSItem

	if leftItem, middleItem, rightItem, err = s.checkoutConflictedForMerge(
		conflicted,
		mode,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var skuReplacement *sku.Transacted
	var replacement *sku.FSItem

	if skuReplacement, err = s.MakeMergedTransacted(conflicted); err != nil {
		var mergeConflict *sku.ErrMergeConflict

		if errors.As(err, &mergeConflict) {
			err = nil
			replacement = &mergeConflict.FSItem
		} else {
			err = errors.Wrap(err)
			return
		}
	} else {
		if replacement, err = s.ReadFSItemFromExternal(skuReplacement); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	tool = append(
		tool,
		leftItem.Object.GetPath(),
		middleItem.Object.GetPath(),
		rightItem.Object.GetPath(),
		replacement.Object.GetPath(),
	)

	// TODO merge blobs

	cmd := exec.Command(tool[0], tool[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	ui.Log().Print(cmd.Env)

	if err = cmd.Run(); err != nil {
		err = errors.Wrapf(err, "Cmd: %q", tool)
		return
	}

	e := GetExternalPool().Get()
	defer GetExternalPool().Put(e)

	e.ObjectId.ResetWith(&co.GetSkuExternal().ObjectId)

	if err = s.WriteFSItemToExternal(leftItem, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = files.Open(replacement.Object.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO open blob

	defer errors.DeferredCloser(&err, f)

	if err = s.ReadOneExternalObjectReader(f, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.DeleteCheckedOut(
		conflicted.CheckedOut,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	co = GetCheckedOutPool().Get()

	sku.TransactedResetter.ResetWith(co.GetSkuExternal(), e)

	return
}
