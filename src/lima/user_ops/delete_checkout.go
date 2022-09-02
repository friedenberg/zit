package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/bravo/stdprinter"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
	"github.com/friedenberg/zit/zettel_external"
	"github.com/friedenberg/zit/zettel_transacted"
)

type DeleteCheckout struct {
	Umwelt *umwelt.Umwelt
}

func (c DeleteCheckout) Run(
	store store_with_lock.Store,
	zettels []zettel_external.Zettel,
) (err error) {
	toDelete := make([]zettel_external.Zettel, 0, len(zettels))
	filesToDelete := make([]string, 0, len(zettels))

	for _, external := range zettels {
		var internal zettel_transacted.Transacted

		if internal, err = store.StoreObjekten().Read(external.Named.Hinweis); err != nil {
			err = errors.Error(err)
			return
		}

		//TODO add a safety check?
		if !internal.Named.Stored.Zettel.Equals(external.Named.Stored.Zettel) {
			stdprinter.Outf("[%s] (checkout different!)\n", external.Named.Hinweis)
			continue
		}

		toDelete = append(toDelete, external)
		filesToDelete = append(filesToDelete, external.ZettelFD.Path)

		if external.AkteFD.Path != "" {
			filesToDelete = append(filesToDelete, external.AkteFD.Path)
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
