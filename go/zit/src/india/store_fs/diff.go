package store_fs

import (
	"io"
	"os"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
)

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

	if f, err = s.fs_home.FileTempLocal(); err != nil {
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
		if cmd.ProcessState.ExitCode() == 1 {
			err = errors.Wrap(MakeErrMergeConflict(nil))
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	path = f.Name()

	return
}
