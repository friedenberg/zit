package zettel

// TODO-P2 deprecate and move to writerfuncs
type Writer interface {
	WriteZettel(z *Zettel) (err error)
}

type WriterFunc func(*Zettel) error

func MakeWriter(f WriterFunc) Writer {
	return WriterFunc(f)
}

func (wf WriterFunc) WriteZettel(z *Zettel) (err error) {
	return wf(z)
}
