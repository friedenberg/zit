package zettel_transacted

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
