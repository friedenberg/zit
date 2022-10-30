package zettel_named

import "github.com/friedenberg/zit/src/bravo/collections"

type Writer interface {
	WriteZettelNamed(*Zettel) (err error)
}

type WriterFunc collections.WriterFunc[*Zettel]

func (wf WriterFunc) WriteZettelNamed(z *Zettel) (err error) {
	return collections.WriterFunc[*Zettel](wf)(z)
}
