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
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

func (s *Store) Merge(tm sku.Conflicted) (err error) {
	p, mergeResult := s.tryMergeIgnoringConflicts(tm)

	var merged FDPair

	if err = merged.Object.SetPath(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	cofs := tm.CheckedOutLike.(*CheckedOut)

	if mergeResult != nil {
		mergeConflict := &ErrMergeConflict{}

		if errors.As(mergeResult, &mergeConflict) {
			if err = s.handleMergeResult(
				tm,
				cofs,
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

	src := merged.Object.GetPath()
	dst := cofs.External.FDs.Object.GetPath()

	// TODO determine why dst is sometimes ""
	if err = files.Rename(src, dst); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) tryMergeIgnoringConflicts(
	tm sku.Conflicted,
) (path string, err error) {
	if err = tm.MergeTags(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var leftCO, middleCO, rightCO *CheckedOut

	inlineBlob := tm.IsAllInlineType(s.config)

	mode := checkout_mode.MetadataAndBlob

	if !inlineBlob {
		mode = checkout_mode.MetadataOnly
	}

	if leftCO, err = s.checkoutOneForMerge(mode, tm.Left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if middleCO, err = s.checkoutOneForMerge(mode, tm.Middle); err != nil {
		err = errors.Wrap(err)
		return
	}

	if rightCO, err = s.checkoutOneForMerge(mode, tm.Right); err != nil {
		err = errors.Wrap(err)
		return
	}

	path, err = s.runDiff3(
		leftCO.External.FDs.Object,
		middleCO.External.FDs.Object,
		rightCO.External.FDs.Object,
	)

	return
}

func (s *Store) checkoutOneForMerge(
	mode checkout_mode.Mode,
	sz *sku.Transacted,
) (cz *CheckedOut, err error) {
	options := checkout_options.Options{
		CheckoutMode:    mode,
		AllowConflicted: true,
		Path:            checkout_options.PathTempLocal,
		ForceInlineBlob: true,
		Force:           true,
	}

	cz = GetCheckedOutPool().Get()
	cz.External.FDs.Reset()
	sku.Resetter.ResetWith(&cz.Internal, sz)

	if err = s.checkoutOne(
		options,
		cz,
	); err != nil {
		ui.Log().Print(&cz.External.FDs.Object, &cz.External.FDs.Blob)
		err = errors.Wrap(err)
		return
	}

	ui.Log().Print(&cz.External.FDs.Object, &cz.External.FDs.Blob)

	return
}

func (s *Store) handleMergeResult(
	conflicted sku.Conflicted,
	cofs *CheckedOut,
	mergeResult *ErrMergeConflict,
) (err error) {
	var f *os.File

	if f, err = s.fs_home.FileTempLocal(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	bw := bufio.NewWriter(f)
	defer errors.DeferredFlusher(&err, bw)

	p := sku_fmt.MakeFormatInventoryListPrinter(
		bw,
		object_inventory_format.FormatForVersion(s.config.GetStoreVersion()),
		s.objectFormatOptions,
	)

	if err = conflicted.WriteConflictMarker(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = cofs.External.FDs.GenerateConflictFD(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Rename(
		f.Name(),
		cofs.External.FDs.Conflict.GetPath(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	cofs.State = checked_out_state.Conflicted

	return
}

func (s *Store) RunMergeTool(
	tool []string,
	tm sku.Conflicted,
) (co *CheckedOut, err error) {
	if len(tool) == 0 {
		err = errors.Errorf("no utility provided")
		return
	}

	co = tm.CheckedOutLike.(*CheckedOut)

	inlineBlob := tm.IsAllInlineType(s.config)

	mode := checkout_mode.MetadataAndBlob

	if !inlineBlob {
		mode = checkout_mode.MetadataOnly
	}

	var leftCO, middleCO, rightCO *CheckedOut

	if leftCO, err = s.checkoutOneForMerge(mode, tm.Left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if middleCO, err = s.checkoutOneForMerge(mode, tm.Middle); err != nil {
		err = errors.Wrap(err)
		return
	}

	if rightCO, err = s.checkoutOneForMerge(mode, tm.Right); err != nil {
		err = errors.Wrap(err)
		return
	}

	tmpPath, mergeResult := s.tryMergeIgnoringConflicts(tm)

	if !errors.Is(mergeResult, &ErrMergeConflict{}) {
		err = errors.Wrap(mergeResult)
		return
	}

	tool = append(tool,
		leftCO.External.FDs.Object.GetPath(),
		middleCO.External.FDs.Object.GetPath(),
		rightCO.External.FDs.Object.GetPath(),
		tmpPath,
	)

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

	e.ResetWith(&leftCO.External)

	var f *os.File

	if f, err = files.Open(tmpPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if err = s.ReadOneExternalObjectReader(f, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.DeleteExternalLike(
		tm.CheckedOutLike.GetSkuExternalLike(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	co = GetCheckedOutPool().Get()

	co.External.ResetWith(e)

	return
}

func MakeErrMergeConflict(sk *FDPair) (err *ErrMergeConflict) {
	err = &ErrMergeConflict{}

	if sk != nil {
		err.ResetWith(sk)
	}

	return
}

type ErrMergeConflict struct {
	FDPair
}

func (e *ErrMergeConflict) Is(target error) bool {
	_, ok := target.(*ErrMergeConflict)
	return ok
}

func (e *ErrMergeConflict) Error() string {
	return fmt.Sprintf(
		"merge conflict for fds: Object: %q, Blob: %q",
		&e.FDPair.Object,
		&e.FDPair.Blob,
	)
}
