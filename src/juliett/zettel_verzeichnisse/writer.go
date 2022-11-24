package zettel_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/hotel/zettel"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
)

func MakeWriterZettelTransacted(
	wf collections.WriterFunc[*zettel_transacted.Zettel],
) collections.WriterFunc[*Zettel] {
	return func(z *Zettel) (err error) {
		if err = wf(&z.Transacted); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func MakeWriterZettelNamed(
	wf collections.WriterFunc[*zettel.Named],
) collections.WriterFunc[*Zettel] {
	return func(z *Zettel) (err error) {
		if err = wf(&z.Transacted.Named); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
