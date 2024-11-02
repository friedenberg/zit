package store_fs

import (
	"io"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
)

func (s *Store) runDiff3(left, middle, right *fd.FD) (path string, err error) {
	cmd := exec.Command(
		"git",
		"merge-file",
		"-p",
		"-L=left",
		"-L=middle",
		"-L=right",
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

	if f, err = s.fs_home.TempLocal.FileTemp(); err != nil {
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

	if err = cmd.Wait(); err != nil {
		var errExit *exec.ExitError

		if !errors.As(err, &errExit) {
			err = errors.Wrap(err)
			return
		}

		// TODO figure out why exit code 2 is being thrown by diff3 for conflicts
		err = errors.Wrap(MakeErrMergeConflict(nil))
		// if errExit.ExitCode() == 1 {
		// } else {
		// 	err = errors.Wrapf(errExit, "Stderr: %q", errExit.Stderr)
		// 	return
		// }
	}

	path = f.Name()

	return
}
