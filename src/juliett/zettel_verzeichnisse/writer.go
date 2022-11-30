package zettel_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/india/zettel"
)

func MakeWriterZettelTransacted(
	wf collections.WriterFunc[*zettel.Transacted],
) collections.WriterFunc[*Zettel] {
	return func(z *Zettel) (err error) {
		if err = wf(&z.Transacted); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
