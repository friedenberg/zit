package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
)

func MakeWriterZettel(
	wf collections.WriterFunc[*Zettel],
) collections.WriterFunc[*Transacted] {
	return func(z *Transacted) (err error) {
		if err = wf(&z.Objekte); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
