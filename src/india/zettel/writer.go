package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
)

func MakeWriterZettelNamed(
	wf collections.WriterFunc[*Named],
) collections.WriterFunc[*Transacted] {
	return func(z *Transacted) (err error) {
		if err = wf(&z.Named); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func MakeWriterZettel(
	wf collections.WriterFunc[*Zettel],
) collections.WriterFunc[*Transacted] {
	return func(z *Transacted) (err error) {
		if err = wf(&z.Named.Stored.Objekte); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
