package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type DeleteCheckout struct {
	*umwelt.Umwelt
}

func (c DeleteCheckout) Run(
	zes zettel_external.MutableSet,
) (err error) {
	zesToDelete := zettel_external.MakeMutableSetUniqueFD()
	filesToDelete := collections.MakeMutableSet[*zettel_external.FD](
		func(e *zettel_external.FD) string {
			if e == nil {
				return ""
			}

			return e.Path
		},
	)

	if err = zes.Each(
		func(external *zettel_external.Zettel) (err error) {
			var internal zettel_transacted.Zettel

			if internal, err = c.StoreObjekten().Zettel().ReadHinweisSchwanzen(external.Named.Kennung); err != nil {
				err = errors.Wrap(err)
				return
			}

			//TODO add a safety check?
			if !internal.Named.Stored.Objekte.Equals(&external.Named.Stored.Objekte) {
				errors.PrintOutf("[%s] (checkout different!)", external.Named.Kennung)
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
		func(e *zettel_external.FD) (err error) {
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
