package user_ops

import (
	"fmt"

	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type DeleteCheckout struct {
}

func (c DeleteCheckout) Run(zettels map[hinweis.Hinweis]stored_zettel.CheckedOut) (err error) {
	toDelete := make([]stored_zettel.External, 0, len(zettels))
	filesToDelete := make([]string, 0, len(zettels))

	for h, z := range zettels {
		if !z.Internal.Zettel.Equals(z.External.Zettel) {
			fmt.Printf("%#v\n", z.Internal.Zettel)
			fmt.Printf("%#v\n", z.External.Zettel)
			_Outf("[%s] (checkout different!)\n", h)
			continue
		}

		toDelete = append(toDelete, z.External)
		filesToDelete = append(filesToDelete, z.External.Path)

		if z.External.AktePath != "" {
			filesToDelete = append(filesToDelete, z.External.AktePath)
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
