package zettel_transacted

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
)

type Writer interface {
	WriteZettelTransacted(*Zettel) error
}

type writer collections.WriterFunc[*Zettel]

func MakeWriter(f collections.WriterFunc[*Zettel]) writer {
	return writer(f)
}

func (w writer) WriteZettelTransacted(z *Zettel) (err error) {
	return collections.WriterFunc[*Zettel](w)(z)
}

func MakeWriterZettel(wf collections.WriterFunc[*zettel.Zettel]) Writer {
	return MakeWriter(
		func(z *Zettel) (err error) {
			if err = wf(&z.Named.Stored.Zettel); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
}

type WriterZettelNamed struct {
	zettel_named.Writer
}

func (w WriterZettelNamed) WriteZettelTransacted(z *Zettel) (err error) {
	return w.WriteZettelNamed(&z.Named)
}
