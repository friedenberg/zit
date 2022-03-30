package user_ops

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type DeleteCheckout struct {
	Umwelt *umwelt.Umwelt
}

func (c DeleteCheckout) Run(zettels map[hinweis.Hinweis]stored_zettel.External) (err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	toDelete := make([]stored_zettel.External, 0, len(zettels))
	filesToDelete := make([]string, 0, len(zettels))

	for h, external := range zettels {
		var internal stored_zettel.Named

		if internal, err = store.Zettels().Read(h); err != nil {
			err = _Error(err)
			return
		}

		//TODO add a safety check?
		if !internal.Zettel.Equals(external.Zettel) {
			_Outf("[%s] (checkout different!)\n", h)
			continue
		}

		toDelete = append(toDelete, external)
		filesToDelete = append(filesToDelete, external.Path)

		if external.AktePath != "" {
			filesToDelete = append(filesToDelete, external.AktePath)
		}
	}

	if err = open_file_guard.DeleteFilesAndDirs(filesToDelete...); err != nil {
		err = _Error(err)
		return
	}

	for _, z := range toDelete {
		_Outf("[%s] (checkout deleted)\n", z.Hinweis)
	}

	return
}
