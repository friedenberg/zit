package store_fs

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
)

type DeleteCheckout struct{}

func (c DeleteCheckout) Run(
	dryRun bool,
	s dir_layout.DirLayout,
	p interfaces.FuncIter[*fd.FD],
	fs interfaces.Iterable[*fd.FD],
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

		if pRel, pErr := filepath.Rel(s.Cwd(), fd.String()); pErr == nil {
			path = pRel
		}

		func() {
			if fd.IsDir() && fd.GetPath() != s.Cwd() {
				dirs = append(dirs, fd.GetPath())
				return
			}

			dir := filepath.Dir(path)

			if dir == s.Cwd() {
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
