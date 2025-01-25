package store_fs

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
)

type DeleteCheckout struct{}

func (c DeleteCheckout) Run(
	dryRun bool,
	s env_repo.Env,
	p interfaces.FuncIter[*fd.FD],
	fs interfaces.Collection[*fd.FD],
) (err error) {
	els := quiter.ElementsSorted(
		fs,
		func(i, j *fd.FD) bool {
			return i.String() < j.String()
		},
	)

	if dryRun {
		for _, f := range els {
			if err = p(f); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return
	}

	dirs := make([]string, 0)

	for _, fd := range els {
		path := fd.String()

		if path == "." {
			continue
		}

		if pRel, pErr := filepath.Rel(
			s.GetCwd(),
			fd.String(),
		); pErr == nil {
			path = pRel
		}

		func() {
			if fd.IsDir() && fd.GetPath() != s.GetCwd() {
				dirs = append(dirs, fd.GetPath())
				return
			}

			dir := filepath.Dir(path)

			if dir == s.GetCwd() {
				return
			}

			// Occurs when the file is top-level relative to Cwd, and so has no parent
			// directory. This is a false-positive and should be ignored.
			if dir == "." {
				return
			}

			dirs = append(dirs, dir)
		}()

		if err = s.Delete(path); err != nil {
			if errors.IsNotExist(err) {
				err = nil
			} else {
				err = errors.Wrapf(err, "FD: %s", fd)
				return
			}
		}

		if p != nil {
			if err = p(fd); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	for _, d := range dirs {
		var contents []os.DirEntry

		if contents, err = files.ReadDir(d); err != nil {
			if errors.IsNotExist(err) {
				err = nil
			} else {
				err = errors.Wrapf(err, "Dir: %s", d)
			}

			continue
		}

		if len(contents) != 0 {
			continue
		}

		if err = s.Delete(d); err != nil {
			err = errors.Wrap(err)
			return
		}

		var f *fd.FD

		if f, err = fd.MakeFromDirPath(d); err != nil {
			err = errors.Wrap(err)
			return
		}

		if p != nil {
			if err = p(f); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}
