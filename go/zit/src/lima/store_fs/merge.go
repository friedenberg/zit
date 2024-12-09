package store_fs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) Merge(tm sku.Conflicted) (err error) {
	var original, replacement *sku.FSItem

	original, replacement, mergeResult := s.tryMergeIgnoringConflicts(tm)

	co := tm.CheckedOut

	if mergeResult != nil {
		mergeConflict := &ErrMergeConflict{}

		if errors.As(mergeResult, &mergeConflict) {
			if err = s.handleMergeResult(
				tm,
				co,
				mergeConflict,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			err = errors.Wrap(mergeResult)
		}

		return
	}

	if !replacement.Object.IsEmpty() {
		if err = files.Rename(
			replacement.Object.GetPath(),
			original.Object.GetPath(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !replacement.Blob.IsEmpty() {
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
) (left, middle, right *sku.FSItem, err error) {
	if _, left, err = s.checkoutOneForMerge(mode, tm.Left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, middle, err = s.checkoutOneForMerge(mode, tm.Middle); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, right, err = s.checkoutOneForMerge(mode, tm.Right); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) tryMergeIgnoringConflicts(
	tm sku.Conflicted,
) (original, replacement *sku.FSItem, err error) {
	if err = tm.MergeTags(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var leftItem, middleItem, rightItem *sku.FSItem

	inlineBlob := tm.IsAllInlineType(s.config)

	mode := checkout_mode.MetadataAndBlob

	if !inlineBlob {
		mode = checkout_mode.MetadataOnly
	}

	if leftItem, middleItem, rightItem, err = s.checkoutConflictedForMerge(
		tm,
		mode,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if original, err = s.ReadFSItemFromExternal(tm.CheckedOut.GetSkuExternal()); err != nil {
		err = errors.Wrap(err)
		return
	}

	{
		var path string

		var diff3Error error

		path, diff3Error = s.runDiff3(
			&leftItem.Object,
			&middleItem.Object,
			&rightItem.Object,
		)

		replacement = &sku.FSItem{}

		if err = replacement.Object.SetPath(path); err != nil {
			err = errors.Wrap(err)
			return
		}

		if diff3Error != nil {
			err = errors.Wrap(diff3Error)
			return
		}
	}

	return
}

func (s *Store) checkoutOneForMerge(
	mode checkout_mode.Mode,
	sz *sku.Transacted,
) (cz *sku.CheckedOut, i *sku.FSItem, err error) {
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

func (s *Store) handleMergeResult(
	conflicted sku.Conflicted,
	cofs *sku.CheckedOut,
	mergeResult *ErrMergeConflict,
) (err error) {
	var f *os.File

	if f, err = s.dirLayout.TempLocal.FileTemp(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	bw := bufio.NewWriter(f)
	defer errors.DeferredFlusher(&err, bw)

	bs := s.externalStoreSupplies.BlobStore.GetInventoryList()

	if _, err = bs.WriteBlobToWriter(
		builtin_types.DefaultOrPanic(genres.InventoryList),
		conflicted,
		bw,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var i *sku.FSItem

	if i, err = s.ReadFSItemFromExternal(cofs.GetSkuExternal()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.GenerateConflictFD(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Rename(
		f.Name(),
		i.Conflict.GetPath(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	cofs.SetState(checked_out_state.Conflicted)

	return
}

func (s *Store) RunMergeTool(
	tool []string,
	tm sku.Conflicted,
) (co *sku.CheckedOut, err error) {
	if len(tool) == 0 {
		err = errors.Errorf("no utility provided")
		return
	}

	co = tm.CheckedOut

	inlineBlob := tm.IsAllInlineType(s.config)

	mode := checkout_mode.MetadataAndBlob

	if !inlineBlob {
		mode = checkout_mode.MetadataOnly
	}

	var leftItem, middleItem, rightItem *sku.FSItem

	if leftItem, middleItem, rightItem, err = s.checkoutConflictedForMerge(
		tm,
		mode,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	_, after, mergeResult := s.tryMergeIgnoringConflicts(tm)

	if !errors.Is(mergeResult, &ErrMergeConflict{}) {
		err = errors.Wrap(mergeResult)
		return
	}

	tool = append(
		tool,
		leftItem.Object.GetPath(),
		middleItem.Object.GetPath(),
		rightItem.Object.GetPath(),
		after.Object.GetPath(),
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

	if f, err = files.Open(after.Object.GetPath()); err != nil {
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
		tm.CheckedOut,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	co = GetCheckedOutPool().Get()

	sku.TransactedResetter.ResetWith(co.GetSkuExternal(), e)

	return
}

func MakeErrMergeConflict(sk *sku.FSItem) (err *ErrMergeConflict) {
	err = &ErrMergeConflict{}

	if sk != nil {
		err.ResetWith(sk)
	}

	return
}

type ErrMergeConflict struct {
	sku.FSItem
}

func (e *ErrMergeConflict) Is(target error) bool {
	_, ok := target.(*ErrMergeConflict)
	return ok
}

func (e *ErrMergeConflict) Error() string {
	return fmt.Sprintf(
		"merge conflict for fds: Object: %q, Blob: %q",
		&e.Object,
		&e.Blob,
	)
}
