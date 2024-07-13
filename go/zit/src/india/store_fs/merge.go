package store_fs

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

func (s *Store) Merge(tm sku.Conflicted) (err error) {
	cofs := tm.CheckedOutLike.(*CheckedOut)

	if err = tm.MergeEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var leftCO, middleCO, rightCO *CheckedOut

	inlineAkte := tm.IsAllInlineTyp(s.konfig)

	mode := checkout_mode.ModeObjekteAndAkte

	if !inlineAkte {
		mode = checkout_mode.ModeObjekteOnly
	}

	op := checkout_options.Options{
		CheckoutMode:    mode,
		ForceInlineAkte: true,
		Path:            checkout_options.PathTempLocal,
		Force:           true,
	}

	if leftCO, err = s.checkoutOneNew(op, tm.Left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if middleCO, err = s.checkoutOneNew(op, tm.Middle); err != nil {
		err = errors.Wrap(err)
		return
	}

	if rightCO, err = s.checkoutOneNew(op, tm.Right); err != nil {
		err = errors.Wrap(err)
		return
	}

	p, mergeResult := s.runDiff3(
		leftCO.External.FDs.Objekte,
		middleCO.External.FDs.Objekte,
		rightCO.External.FDs.Objekte,
	)

	var merged FDPair

	if err = merged.Objekte.SetPath(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mergeResult != nil {
		mergeConflict := &ErrMergeConflict{}

		if errors.As(mergeResult, &mergeConflict) {
			if err = s.handleMergeResult(
				tm,
				cofs,
				mergeConflict,
				merged,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			err = errors.Wrap(mergeResult)
		}

		return
	}

	// if merged.Akte.Path, err = s.runDiff3(
	// 	leftCO.External.FDs.Akte,
	// 	middleCO.External.FDs.Akte,
	// 	rightCO.External.FDs.Akte,
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }
	src := merged.Objekte.GetPath()
	dst := cofs.External.FDs.Objekte.GetPath()

	if err = files.Rename(src, dst); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) handleMergeResult(
	conflicted sku.Conflicted,
	cofs *CheckedOut,
	mergeResult *ErrMergeConflict,
	merged FDPair,
) (err error) {
	var f *os.File

	if f, err = s.fs_home.FileTempLocal(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	bw := bufio.NewWriter(f)
	defer errors.DeferredFlusher(&err, bw)

	p := sku_fmt.MakeFormatBestandsaufnahmePrinter(
		bw,
		objekte_format.FormatForVersion(s.konfig.GetStoreVersion()),
		s.objekteFormatOptions,
	)

	if err = conflicted.WriteConflictMarker(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Rename(
		f.Name(),
		cofs.External.FDs.MakeConflictMarker(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) RunMergeTool(
	tool []string,
	tm sku.Conflicted,
) (co *CheckedOut, err error) {
	inlineAkte := tm.IsAllInlineTyp(s.konfig)

	op := checkout_options.Options{
		CheckoutMode:    checkout_mode.ModeObjekteAndAkte,
		AllowConflicted: true,
		Path:            checkout_options.PathTempLocal,
	}

	if !inlineAkte {
		op.CheckoutMode = checkout_mode.ModeObjekteOnly
	}

	var leftCO, middleCO, rightCO *CheckedOut

	if leftCO, err = s.checkoutOneNew(op, tm.Left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if middleCO, err = s.checkoutOneNew(op, tm.Middle); err != nil {
		err = errors.Wrap(err)
		return
	}

	if rightCO, err = s.checkoutOneNew(op, tm.Right); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = s.fs_home.FileTempLocal(); err != nil {
		err = errors.Wrap(err)
		return
	}

	tmpPath := f.Name()

	defer errors.DeferredCloser(&err, f)

	var cmdStrings []string

	if len(tool) > 1 {
		toolArgs := tool[1:]
		cmdStrings = make([]string, 3+len(toolArgs))
		copy(cmdStrings, toolArgs)
	}

	cmdStrings = append(
		cmdStrings,
		leftCO.External.FDs.Objekte.GetPath(),
		middleCO.External.FDs.Objekte.GetPath(),
		rightCO.External.FDs.Objekte.GetPath(),
		tmpPath,
	)

	if len(tool) == 0 {
		err = errors.Errorf("no utility provided")
		return
	}

	cmd := exec.Command(tool[0], cmdStrings...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{
		fmt.Sprintf("LOCAL=%s", leftCO.External.FDs.Objekte.GetPath()),
		fmt.Sprintf("BASE=%s", middleCO.External.FDs.Objekte.GetPath()),
		fmt.Sprintf("REMOTE=%s", rightCO.External.FDs.Objekte.GetPath()),
		fmt.Sprintf("MERGED=%s", tmpPath),
	}

	if err = cmd.Run(); err != nil {
		err = errors.Wrap(err)
		return
	}

	e := GetExternalPool().Get()
	defer GetExternalPool().Put(e)

	if err = e.SetFromSkuLike(&leftCO.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadOneExternalObjekteReader(f, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.DeleteCheckout(tm.CheckedOutLike); err != nil {
		err = errors.Wrap(err)
		return
	}

	co = GetCheckedOutPool().Get()

	if err = co.External.SetFromSkuLike(e); err != nil {
		err = errors.Wrap(err)
		return
	}

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
	return fmt.Sprintf("merge conflict for fds: %v", &e.FDPair)
}
