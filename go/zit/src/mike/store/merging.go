package store

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/to_merge"
)

func (s *Store) readExternalAndMergeIfNecessary(
	transactedPtr, mutter *sku.Transacted,
) (err error) {
	if mutter == nil {
		return
	}

	var co *store_fs.CheckedOut

	if co, err = s.CombineOneCheckedOutFS(mutter); err != nil {
		err = nil
		return
	}

	defer store_fs.GetCheckedOutPool().Put(co)

	mutterEqualsExternal := co.InternalAndExternalEqualsSansTai()

	if mutterEqualsExternal {
		var mode checkout_mode.Mode

		if mode, err = co.External.GetCheckoutMode(); err != nil {
			err = errors.Wrap(err)
			return
		}

		op := checkout_options.Options{
			CheckoutMode: mode,
			Force:        true,
		}

		if co, err = s.CheckoutOneFS(op, transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		store_fs.GetCheckedOutPool().Put(co)

		return
	}

	transactedPtrCopy := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(transactedPtrCopy)

	if err = transactedPtrCopy.SetFromSkuLike(transactedPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	tm := to_merge.Sku{
		Left:   transactedPtrCopy,
		Middle: &co.Internal,
		Right:  &co.External.Transacted,
	}

	var merged store_fs.FDPair

	merged, err = s.merge(tm)

	switch {
	case errors.Is(err, &to_merge.ErrMergeConflict{}):
		if err = tm.WriteConflictMarker(
			s.GetStandort(),
			s.GetKonfig().GetStoreVersion(),
			s.GetObjekteFormatOptions(),
			co.External.FDs.MakeConflictMarker(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case err != nil:
		err = errors.Wrap(err)
		return

	default:
		src := merged.Objekte.GetPath()
		dst := co.External.FDs.Objekte.GetPath()

		if err = files.Rename(src, dst); err != nil {
			return
		}
	}

	return
}

func (s *Store) merge(tm to_merge.Sku) (merged store_fs.FDPair, err error) {
	if err = tm.MergeEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var leftCO, middleCO, rightCO *store_fs.CheckedOut

	inlineAkte := tm.IsAllInlineTyp(s.GetKonfig())

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

	if leftCO, err = s.CheckoutOneFS(op, tm.Left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if middleCO, err = s.CheckoutOneFS(op, tm.Middle); err != nil {
		err = errors.Wrap(err)
		return
	}

	if rightCO, err = s.CheckoutOneFS(op, tm.Right); err != nil {
		err = errors.Wrap(err)
		return
	}

	var p string

	if p, err = s.runDiff3(
		leftCO.External.FDs.Objekte,
		middleCO.External.FDs.Objekte,
		rightCO.External.FDs.Objekte,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = merged.Objekte.SetPath(p); err != nil {
		err = errors.Wrap(err)
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

	return
}

func (s *Store) runDiff3(left, middle, right fd.FD) (path string, err error) {
	cmd := exec.Command(
		"diff3",
		"--text",
		"--merge",
		"--label=left",
		"--label=middle",
		"--label=right",
		left.GetPath(),
		middle.GetPath(),
		right.GetPath(),
	)

	var out io.ReadCloser

	if out, err = cmd.StdoutPipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

	cmd.Stderr = os.Stderr

	var f *os.File

	if f, err = s.GetStandort().FileTempLocal(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if err = cmd.Start(); err != nil {
		err = errors.Wrap(err)
		return
	}

	// merged = sku.GetExternalPool().Get()

	// if err = merged.Kennung.SetWithKennung(tm.right.Kennung); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	if _, err = io.Copy(f, out); err != nil {
		err = errors.Wrap(err)
		return
	}

	// if err = s.ReadOneExternalObjekteReader(out, merged); err != nil {
	// 	log.Debug().Printf("%s", err)
	// 	err = nil
	// 	// err = errors.Wrap(err)
	// 	// return
	// }

	if err = cmd.Wait(); err != nil {
		if cmd.ProcessState.ExitCode() == 1 {
			err = errors.Wrap(to_merge.MakeErrMergeConflict(nil))
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	path = f.Name()

	return
}

func (s *Store) RunMergeTool(
	tm to_merge.Sku,
) (err error) {
	tool := s.GetKonfig().Cli().ToolOptions.Merge
	inlineAkte := tm.IsAllInlineTyp(s.GetKonfig())

	op := checkout_options.Options{
		CheckoutMode:    checkout_mode.ModeObjekteAndAkte,
		AllowConflicted: true,
		Path:            checkout_options.PathTempLocal,
	}

	if !inlineAkte {
		op.CheckoutMode = checkout_mode.ModeObjekteOnly
	}

	var leftCO, middleCO, rightCO *store_fs.CheckedOut

	if leftCO, err = s.CheckoutOneFS(op, tm.Left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if middleCO, err = s.CheckoutOneFS(op, tm.Middle); err != nil {
		err = errors.Wrap(err)
		return
	}

	if rightCO, err = s.CheckoutOneFS(op, tm.Right); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = s.GetStandort().FileTempLocal(); err != nil {
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

	e := store_fs.GetExternalPool().Get()
	defer store_fs.GetExternalPool().Put(e)

	if err = e.SetFromSkuLike(&leftCO.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadOneExternalObjekteReader(f, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetStandort().Delete(tm.ConflictMarkerPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	co := store_fs.GetCheckedOutPool().Get()
	defer store_fs.GetCheckedOutPool().Put(co)

	if err = co.External.SetFromSkuLike(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = s.CreateOrUpdateCheckedOutFS(co, false); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
