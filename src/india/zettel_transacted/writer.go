package zettel_transacted

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/hotel/zettel"
)

func MakeWriterZettelNamed(
	wf collections.WriterFunc[*zettel.Named],
) collections.WriterFunc[*Zettel] {
	return func(z *Zettel) (err error) {
		if err = wf(&z.Named); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func MakeWriterZettel(
	wf collections.WriterFunc[*zettel.Zettel],
) collections.WriterFunc[*Zettel] {
	return func(z *Zettel) (err error) {
		if err = wf(&z.Named.Stored.Objekte); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
