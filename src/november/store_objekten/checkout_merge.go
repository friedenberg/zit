package store_objekten

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/charlie/checkout_options"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/to_merge"
)

func (s *Store) merge(tm to_merge.Sku) (merged sku.ExternalFDs, err error) {
	if err = tm.MergeEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var leftCO, middleCO, rightCO *sku.CheckedOut

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

	if leftCO, err = s.CheckoutOne(op, tm.Left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if middleCO, err = s.CheckoutOne(op, tm.Middle); err != nil {
		err = errors.Wrap(err)
		return
	}

	if rightCO, err = s.CheckoutOne(op, tm.Right); err != nil {
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
			err = errors.Wrap(to_merge.ErrMergeConflict{})
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
	}

	if !inlineAkte {
		op.CheckoutMode = checkout_mode.ModeObjekteOnly
	}

	var leftCO, middleCO, rightCO *sku.CheckedOut

	if leftCO, err = s.CheckoutOne(op, tm.Left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if middleCO, err = s.CheckoutOne(op, tm.Middle); err != nil {
		err = errors.Wrap(err)
		return
	}

	if rightCO, err = s.CheckoutOne(op, tm.Right); err != nil {
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

	e := sku.GetExternalPool().Get()
	defer sku.GetExternalPool().Put(e)

	*e = leftCO.External

	if err = s.ReadOneExternalObjekteReader(f, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Remove(tm.ConflictMarkerPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	co := sku.GetCheckedOutPool().Get()
	defer sku.GetCheckedOutPool().Put(co)

	co.External = *e

	if _, err = s.CreateOrUpdateCheckedOut(co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}