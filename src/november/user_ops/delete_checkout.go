package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/golf/zettel_external"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/mike/umwelt"
)

type DeleteCheckout struct {
	*umwelt.Umwelt
}

func (c DeleteCheckout) Run(
	zettels []zettel_external.Zettel,
) (err error) {
	toDelete := make([]zettel_external.Zettel, 0, len(zettels))
	filesToDelete := make([]string, 0, len(zettels))

	for _, external := range zettels {
		var internal zettel_transacted.Zettel

		if internal, err = c.StoreObjekten().ReadHinweisSchwanzen(external.Named.Hinweis); err != nil {
			err = errors.Wrap(err)
			return
		}

		//TODO add a safety check?
		if !internal.Named.Stored.Zettel.Equals(external.Named.Stored.Zettel) {
			errors.PrintOutf("[%s] (checkout different!)", external.Named.Hinweis)
			continue
		}

		toDelete = append(toDelete, external)
		filesToDelete = append(filesToDelete, external.ZettelFD.Path)

		if external.AkteFD.Path != "" {
			filesToDelete = append(filesToDelete, external.AkteFD.Path)
		}
	}

	if err = files.DeleteFilesAndDirs(filesToDelete...); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, f := range filesToDelete {
		errors.PrintOutf("%s (checkout deleted)", f)
	}

	return
}
