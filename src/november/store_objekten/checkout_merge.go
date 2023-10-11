package store_objekten

import (
	"io"
	"os"
	"os/exec"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type toMerge struct {
	left, middle, right *sku.Transacted
}

func (s *Store) mergeEtiketten(tm toMerge) (err error) {
	left := tm.left.GetEtiketten().CloneMutableSetPtrLike()
	middle := tm.middle.GetEtiketten().CloneMutableSetPtrLike()
	right := tm.right.GetEtiketten().CloneMutableSetPtrLike()

	same := kennung.MakeEtikettMutableSet()
	deleted := kennung.MakeEtikettMutableSet()

	removeFromAllButAddTo := func(
		e *kennung.Etikett,
		toAdd kennung.EtikettMutableSet,
	) (err error) {
		if err = toAdd.AddPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = left.DelPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = middle.DelPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = right.DelPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = middle.EachPtr(
		func(e *kennung.Etikett) (err error) {
			if left.ContainsKey(left.KeyPtr(e)) && right.ContainsKey(right.KeyPtr(e)) {
				return removeFromAllButAddTo(e, same)
			} else if left.ContainsKey(left.KeyPtr(e)) || right.ContainsKey(right.KeyPtr(e)) {
				return removeFromAllButAddTo(e, deleted)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = left.EachPtr(same.AddPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = right.EachPtr(same.AddPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	ets := same.CloneSetPtrLike()

	tm.left.GetMetadateiPtr().Etiketten = ets
	tm.middle.GetMetadateiPtr().Etiketten = ets
	tm.right.GetMetadateiPtr().Etiketten = ets

	return
}

func (s *Store) merge(tm toMerge) (merged sku.ExternalFDs, err error) {
	if err = s.mergeEtiketten(tm); err != nil {
		err = errors.Wrap(err)
		return
	}

	var leftCO, middleCO, rightCO *sku.CheckedOut

	op := store_util.CheckoutOptions{
		CheckoutMode:    checkout_mode.ModeObjekteAndAkte,
		ForceInlineAkte: true,
		UseTempDir:      true,
		Force:           true,
	}

	if leftCO, err = s.CheckoutOne(op, tm.left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if middleCO, err = s.CheckoutOne(op, tm.middle); err != nil {
		err = errors.Wrap(err)
		return
	}

	if rightCO, err = s.CheckoutOne(op, tm.right); err != nil {
		err = errors.Wrap(err)
		return
	}

	if merged.Objekte.Path, err = s.runDiff3(
		leftCO.External.FDs.Objekte,
		middleCO.External.FDs.Objekte,
		rightCO.External.FDs.Objekte,
	); err != nil {
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

func (s *Store) runDiff3(left, middle, right kennung.FD) (path string, err error) {
	cmd := exec.Command(
		"diff3",
		"--text",
		"--merge",
		"--label=left",
		"--label=middle",
		"--label=right",
		left.Path,
		middle.Path,
		right.Path,
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
			err = ErrMergeConflict{
				ExternalFDs: sku.ExternalFDs{
					Objekte: kennung.FD{
						Path: f.Name(),
					},
				},
			}
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	path = f.Name()

	return
}
