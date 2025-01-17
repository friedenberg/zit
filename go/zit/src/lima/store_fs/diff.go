package store_fs

import (
	"io"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

// TODO include blobs
func (s *Store) runDiff3(
	local, base, remote *sku.FSItem,
) (merged *sku.FSItem, err error) {
	baseObjectPath := "/dev/null"

	if base.Len() > 0 {
		baseObjectPath = base.Object.GetPath()
	}

	cmd := exec.Command(
		"git",
		"merge-file",
		"-p",
		"-L=local",
		"-L=base",
		"-L=remote",
		local.Object.GetPath(),
		baseObjectPath,
		remote.Object.GetPath(),
	)

	var out io.ReadCloser

	if out, err = cmd.StdoutPipe(); err != nil {
		err = errors.Wrap(err)
		return
	}

	cmd.Stderr = os.Stderr

	var f *os.File

	if f, err = s.repoLayout.GetDirLayout().TempLocal.FileTemp(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if err = cmd.Start(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = io.Copy(f, out); err != nil {
		err = errors.Wrap(err)
		return
	}

	merged = &sku.FSItem{}
	merged.ResetWith(local)
	merged.Object.Reset()
	merged.Blob.Reset()
	merged.MutableSetLike.Reset()

	hasConflict := false

	if err = cmd.Wait(); err != nil {
		var errExit *exec.ExitError

		if !errors.As(err, &errExit) {
			err = errors.Wrap(err)
			return
		}

		hasConflict = true
	}

	if err = merged.Object.SetPath(f.Name()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if hasConflict {
		err = errors.Wrap(sku.MakeErrMergeConflict(merged))
	}

	return
}
