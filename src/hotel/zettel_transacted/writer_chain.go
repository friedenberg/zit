package zettel_transacted

import "github.com/friedenberg/zit/src/alfa/errors"

type writerChain struct {
	writers []Writer
}

func MakeWriterChain(ws ...Writer) *writerChain {
	return &writerChain{
		writers: ws,
	}
}

func (wc writerChain) WriteZettelTransacted(z *Zettel) (err error) {
	for _, w := range wc.writers {
		err = w.WriteZettelTransacted(z)

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
