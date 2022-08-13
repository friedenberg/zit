package user_ops

import (
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type DeleteCheckout struct {
	Umwelt *umwelt.Umwelt
}

func (c DeleteCheckout) Run(
	store store_with_lock.Store,
	zettels map[hinweis.Hinweis]stored_zettel.External,
) (err error) {
	toDelete := make([]stored_zettel.External, 0, len(zettels))
	filesToDelete := make([]string, 0, len(zettels))

	for h, external := range zettels {
		var internal stored_zettel.Transacted

		if internal, err = store.Zettels().Read(h); err != nil {
			err = errors.Error(err)
			return
		}

		//TODO add a safety check?
		if !internal.Zettel.Equals(external.Zettel) {
			stdprinter.Outf("[%s] (checkout different!)\n", h)
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
