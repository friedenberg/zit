package zettel_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
)

type writerChain struct {
	writers []Writer
}

func MakeWriterChain(ws ...Writer) *writerChain {
	return &writerChain{
		writers: ws,
	}
}

func (wc writerChain) ToPooledChain() collections.PooledChain[Zettel] {
	out := make([]collections.WriterFunc[*Zettel], len(wc.writers))

	for i, wf := range wc.writers {
		out[i] = wf.WriteZettelVerzeichnisse
	}

	return collections.MakePooledChain[Zettel](out...)
}

func (wc writerChain) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	for _, w := range wc.writers {
		err = w.WriteZettelVerzeichnisse(z)

		switch {
		case z == nil && err == nil:
			return

		case err == nil:
			continue

		case errors.IsEOF(err):
			err = nil
			return

		default:
			err = errors.Wrap(err)
			return
		}
	}

	return
}
