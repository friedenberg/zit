package user_ops

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/values"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/charlie/files"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
)

type DeleteCheckout struct {
	*umwelt.Umwelt
}

func (c DeleteCheckout) Run(
	fs schnittstellen.Iterable[*fd.FD],
) (err error) {
	p := c.PrinterFDDeleted()

	if c.Konfig().DryRun {
		return fs.Each(p)
	}

	dirs := collections_value.MakeMutableValueSet[values.String](nil)

	if err = fs.Each(
		func(fd *fd.FD) (err error) {
			path := fd.String()

			if path == "." {
				return
			}

			if pRel, pErr := filepath.Rel(c.Standort().Cwd(), fd.String()); pErr == nil {
				path = pRel
			}

			func() {
				if fd.IsDir() && fd.GetPath() != c.Standort().Cwd() {
					dirs.Add(values.MakeString(fd.GetPath()))
					return
				}

				dir := filepath.Dir(path)

				if dir == c.Standort().Cwd() {
					return
				}

				dirs.Add(values.MakeString(dir))
			}()

			if err = os.Remove(path); err != nil {
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

			if err = os.Remove(d.String()); err != nil {
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
