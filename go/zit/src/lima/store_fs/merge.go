package store_fs

import (
	"bufio"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

// TODO combine with other method in this file
// Makes hard assumptions about the availability of the blobs associated with
// the *sku.CheckedOut.
func (s *Store) MergeCheckedOut(
	co *sku.CheckedOut,
	parentNegotiator sku.ParentNegotiator,
	allowMergeConflicts bool,
) (commitOptions sku.CommitOptions, err error) {
	commitOptions.StoreOptions = sku.GetStoreOptionsImport()

	if co.GetSku().Metadata.Sha().IsNull() || allowMergeConflicts {
		return
	}

	var conflicts checkout_mode.Mode

	// TODO add checkout_mode.BlobOnly
	if co.GetSku().Metadata.Sha().Equals(co.GetSkuExternal().Metadata.Sha()) {
		commitOptions.StoreOptions = sku.StoreOptions{}
		return
	} else if co.GetSku().Metadata.EqualsSansTai(&co.GetSkuExternal().Metadata) {
		if !co.GetSku().Metadata.Tai.Less(co.GetSkuExternal().Metadata.Tai) {
			// TODO implement retroactive change
		}

		return
	} else if co.GetSku().Metadata.Blob.Equals(&co.GetSkuExternal().Metadata.Blob) {
		conflicts = checkout_mode.MetadataOnly
	} else {
		conflicts = checkout_mode.MetadataAndBlob
	}

	// TODO write conflicts
	switch conflicts {
	case checkout_mode.BlobOnly:
	case checkout_mode.MetadataOnly:
	case checkout_mode.MetadataAndBlob:
	default:
	}

	conflicted := sku.Conflicted{
		CheckedOut: co,
		Local:      co.GetSku(),
		Remote:     co.GetSkuExternal(),
	}

	if err = conflicted.FindBestCommonAncestor(parentNegotiator); err != nil {
		err = errors.Wrap(err)
		return
	}

	var skuReplacement *sku.Transacted

	// TODO pass mode / conflicts
	if skuReplacement, err = s.MakeMergedTransacted(
		conflicted,
	); err != nil {
		if sku.IsErrMergeConflict(err) {
			err = nil

			if !allowMergeConflicts {
				if err = s.GenerateConflictMarker(
					conflicted,
					conflicted.CheckedOut,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			co.SetState(checked_out_state.Conflicted)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	sku.TransactedResetter.ResetWith(co.GetSkuExternal(), skuReplacement)

	return
}

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

	localItem.ResetWith(mergedItem)

	merged = GetExternalPool().Get()

	merged.ObjectId.ResetWith(&conflicted.GetSku().ObjectId)

	if err = s.WriteFSItemToExternal(localItem, merged); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.HydrateExternalFromItem(
		sku.CommitOptions{
			StoreOptions: sku.StoreOptions{
				UpdateTai: true,
			},
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
	sk *sku.Transacted,
) (co *sku.CheckedOut, i *sku.FSItem, err error) {
	if sk == nil {
		i = &sku.FSItem{}
		i.Reset()
		return
	}

	options := checkout_options.Options{
		CheckoutMode: mode,
		OptionsWithoutMode: checkout_options.OptionsWithoutMode{
			Force: true,
			StoreSpecificOptions: CheckoutOptions{
				AllowConflicted: true,
				Path:            PathOptionTempLocal,
				// TODO handle binary blobs
				ForceInlineBlob: true,
			},
		},
	}

	co = GetCheckedOutPool().Get()
	sku.Resetter.ResetWith(co.GetSku(), sk)

	if i, err = s.ReadFSItemFromExternal(co.GetSku()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.checkoutOneForReal(
		options,
		co,
		i,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.WriteFSItemToExternal(i, co.GetSkuExternal()); err != nil {
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

	if f, err = s.envRepo.GetTempLocal().FileTemp(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	bw := bufio.NewWriter(f)
	defer errors.DeferredFlusher(&err, bw)

	blobStore := s.storeSupplies.BlobStore.InventoryList

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
			s.envRepo.GetCwd(),
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

func (store *Store) RunMergeTool(
	tool []string,
	conflicted sku.Conflicted,
) (co *sku.CheckedOut, err error) {
	if len(tool) == 0 {
		err = errors.ErrorWithStackf("no utility provided")
		return
	}

	co = conflicted.CheckedOut

	inlineBlob := conflicted.IsAllInlineType(store.config)

	mode := checkout_mode.MetadataAndBlob

	if !inlineBlob {
		mode = checkout_mode.MetadataOnly
	}

	var localItem, baseItem, remoteItem *sku.FSItem

	if localItem, baseItem, remoteItem, err = store.checkoutConflictedForMerge(
		conflicted,
		mode,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var skuReplacement *sku.Transacted
	var replacement *sku.FSItem

	if skuReplacement, err = store.MakeMergedTransacted(conflicted); err != nil {
		var mergeConflict *sku.ErrMergeConflict

		if errors.As(err, &mergeConflict) {
			err = nil
			replacement = &mergeConflict.FSItem
		} else {
			err = errors.Wrap(err)
			return
		}
	} else {
		if replacement, err = store.ReadFSItemFromExternal(skuReplacement); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	tool = append(
		tool,
		localItem.Object.GetPath(),
		baseItem.Object.GetPath(),
		remoteItem.Object.GetPath(),
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

	external := GetExternalPool().Get()
	defer GetExternalPool().Put(external)

	external.ObjectId.ResetWith(&co.GetSkuExternal().ObjectId)

	if err = store.WriteFSItemToExternal(localItem, external); err != nil {
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

	if err = store.ReadOneExternalObjectReader(f, external); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = store.DeleteCheckedOut(
		conflicted.CheckedOut,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	co = GetCheckedOutPool().Get()

	sku.TransactedResetter.ResetWith(co.GetSkuExternal(), external)

	return
}
