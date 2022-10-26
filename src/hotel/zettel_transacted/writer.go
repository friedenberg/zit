package zettel_transacted

import "github.com/friedenberg/zit/src/foxtrot/zettel_named"

type Writer interface {
	WriteZettelTransacted(*Zettel) error
}

type WriterFunc func(*Zettel) error

type writer WriterFunc

func MakeWriter(f WriterFunc) writer {
	return writer(f)
}

func (w writer) WriteZettelTransacted(z *Zettel) (err error) {
	return WriterFunc(w)(z)
}

type WriterZettelNamed struct {
	zettel_named.Writer
}

func (w WriterZettelNamed) WriteZettelTransacted(z *Zettel) (err error) {
	return w.WriteZettelNamed(&z.Named)
}
