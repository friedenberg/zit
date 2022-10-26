package zettel_verzeichnisse

import "github.com/friedenberg/zit/src/alfa/errors"

type writerChain struct {
	writers   []Writer
	ignoreEOF bool
}

func MakeWriterChain(ws ...Writer) *writerChain {
	return &writerChain{
		writers: ws,
	}
}

func MakeWriterChainIgnoreEOF(ws ...Writer) (wc *writerChain) {
	wc = MakeWriterChain(ws...)
	wc.ignoreEOF = true

	return
}

func (wc writerChain) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	for _, w := range wc.writers {
		err = w.WriteZettelVerzeichnisse(z)

		switch {
		case z == nil && err == nil:
			return

		case errors.IsEOF(err) && !wc.ignoreEOF:
			err = nil
			return

		case errors.IsEOF(err) && wc.ignoreEOF:
			err = nil
			continue

		case err == nil:
			continue

		default:
			err = errors.Wrap(err)
			return
		}
	}

	return
}
