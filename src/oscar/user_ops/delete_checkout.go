package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type DeleteCheckout struct {
	*umwelt.Umwelt
}

func (c DeleteCheckout) Run(
	zes zettel_external.MutableSet,
) (err error) {
	zesToDelete := zettel_external.MakeMutableSetUniqueFD()
	filesToDelete := collections.MakeMutableSet[*fd.FD](
		func(e *fd.FD) string {
			if e == nil {
				return ""
			}

			return e.Path
		},
	)

	if err = zes.Each(
		func(external *zettel_external.Zettel) (err error) {
			var internal zettel.Transacted

			if internal, err = c.StoreObjekten().Zettel().ReadHinweisSchwanzen(
				external.Sku.Kennung,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			//TODO add a safety check?
			if !internal.Objekte.Equals(&external.Objekte) {
				errors.Out().Printf("[%s] (checkout different!)", external.Sku.Kennung)
				return
			}

			zesToDelete.Add(external)
			filesToDelete.Add(&external.ZettelFD)

			if external.AkteFD.Path != "" {
				filesToDelete.Add(&external.AkteFD)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	fs := make([]string, 0, filesToDelete.Len())

	filesToDelete.Each(
		func(e *fd.FD) (err error) {
			fs = append(fs, e.Path)
			return
		},
	)

	if err = files.DeleteFilesAndDirs(fs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	filesToDelete.Each(c.PrinterFDDeleted())

	return
}
