package zettel_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
)

type Writer interface {
	WriteZettelVerzeichnisse(z *Zettel) (err error)
}

func MakeWriterZettelTransacted(wf collections.WriterFunc[*zettel_transacted.Zettel]) Writer {
	return MakeWriter(
		func(z *Zettel) (err error) {
			if err = wf(&z.Transacted); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
}

func MakeWriterZettelNamed(wf collections.WriterFunc[*zettel_named.Zettel]) Writer {
	return MakeWriter(
		func(z *Zettel) (err error) {
			if err = wf(&z.Transacted.Named); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
}

type writerFunc collections.WriterFunc[*Zettel]

func (w writerFunc) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	return collections.WriterFunc[*Zettel](w)(z)
}

func MakeWriter(f collections.WriterFunc[*Zettel]) Writer {
	return writerFunc(f)
}
