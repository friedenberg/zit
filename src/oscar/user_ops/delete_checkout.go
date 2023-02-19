package user_ops

import (
	"os"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
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

	return fs.Each(
		func(fd kennung.FD) (err error) {
			path := fd.String()

			if pRel, pErr := filepath.Rel(c.Standort().Cwd(), fd.String()); pErr == nil {
				path = pRel
			}

			if err = os.Remove(path); err != nil {
				if errors.IsNotExist(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			return p(&fd)
		},
	)
}
