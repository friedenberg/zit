package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	"github.com/friedenberg/zit/src/india/store_with_lock"
)

type DeleteCheckout struct {
	Umwelt *umwelt.Umwelt
}

func (c DeleteCheckout) Run(
	store store_with_lock.Store,
	zettels []stored_zettel.External,
) (err error) {
	toDelete := make([]stored_zettel.External, 0, len(zettels))
	filesToDelete := make([]string, 0, len(zettels))

	for _, external := range zettels {
		var internal stored_zettel.Transacted

		if internal, err = store.StoreObjekten().Read(external.Hinweis); err != nil {
			err = errors.Error(err)
			return
		}

		//TODO add a safety check?
		if !internal.Named.Stored.Zettel.Equals(external.Named.Stored.Zettel) {
			stdprinter.Outf("[%s] (checkout different!)\n", external.Named.Hinweis)
			continue
		}

		toDelete = append(toDelete, external)
		filesToDelete = append(filesToDelete, external.Path)

		if external.AktePath != "" {
			filesToDelete = append(filesToDelete, external.AktePath)
		}
	}

	if err = open_file_guard.DeleteFilesAndDirs(filesToDelete...); err != nil {
		err = errors.Error(err)
		return
	}

	for _, f := range filesToDelete {
		stdprinter.Outf("%s (checkout deleted)\n", f)
	}

	return
}
