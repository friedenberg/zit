package store_fs

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
)

type DeleteCheckout struct{}

func (c DeleteCheckout) Run(
	dryRun bool,
	s standort.Standort,
	p schnittstellen.FuncIter[*fd.FD],
	fs schnittstellen.Iterable[*fd.FD],
) (err error) {
	if dryRun {
		if err = fs.Each(p); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	dirs := collections_value.MakeMutableValueSet[values.String](nil)

	if err = fs.Each(
		func(fd *fd.FD) (err error) {
			path := fd.String()

			if path == "." {
				return
			}

			if pRel, pErr := filepath.Rel(s.Cwd(), fd.String()); pErr == nil {
				path = pRel
			}

			func() {
				if fd.IsDir() && fd.GetPath() != s.Cwd() {
					dirs.Add(values.MakeString(fd.GetPath()))
					return
				}

				dir := filepath.Dir(path)

				if dir == s.Cwd() {
					return
				}

				dirs.Add(values.MakeString(dir))
			}()

			if err = s.Delete(path); err != nil {
				if errors.IsNotExist(err) {
					err = nil
				} else {
					err = errors.Wrapf(err, "FD: %s", fd)
					return
				}
			}

			return p(fd)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = dirs.Each(
		func(d values.String) (err error) {
			var contents []os.DirEntry

			if contents, err = files.ReadDir(d.String()); err != nil {
				if errors.IsNotExist(err) {
					err = nil
				} else {
					err = errors.Wrapf(err, "Dir: %s", d)
				}

				return
			}

			if len(contents) != 0 {
				return
			}

			if err = s.Delete(d.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			var f *fd.FD

			if f, err = fd.FDFromDir(d.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return p(f)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
