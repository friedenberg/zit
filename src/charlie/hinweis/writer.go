package hinweis

type Writer interface {
	WriteHinweis(*Hinweis) error
}

type WriterFunc func(*Hinweis) error

type writer WriterFunc

func MakeWriter(f WriterFunc) writer {
	return writer(f)
}

func (w writer) WriteHinweis(h *Hinweis) (err error) {
	return WriterFunc(w)(h)
}
