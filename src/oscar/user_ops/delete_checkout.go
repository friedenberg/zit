package user_ops

import (
	"os"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type DeleteCheckout struct {
	*umwelt.Umwelt
}

func (c DeleteCheckout) Run(
	fs schnittstellen.Set[kennung.FD],
) (err error) {
	p := c.PrinterFDDeleted()

	if c.Konfig().DryRun {
		return fs.EachPtr(p)
	}

	dirs := collections.MakeMutableSetStringer[values.String]()

	if err = fs.Each(
		func(fd kennung.FD) (err error) {
			path := fd.String()

			if path == "." {
				return
			}

			if pRel, pErr := filepath.Rel(c.Standort().Cwd(), fd.String()); pErr == nil {
				path = pRel
			}

			func() {
				if fd.IsDir && fd.Path != c.Standort().Cwd() {
					dirs.Add(values.MakeString(fd.Path))
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

			return p(&fd)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = dirs.Each(
		func(d values.String) (err error) {
			var contents []os.DirEntry

			if contents, err = files.ReadDir(d.String()); err != nil {
				err = errors.Wrapf(err, "Dir: %s", d)
				return
			}

			if len(contents) != 0 {
				return
			}

			if err = os.Remove(d.String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return p(&kennung.FD{Path: d.String(), IsDir: true})
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
