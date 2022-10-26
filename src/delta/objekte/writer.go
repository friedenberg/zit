package objekte

type Writer interface {
	WriteObjekte(Objekte) error
}

type WriterFunc func(Objekte) error

type writer WriterFunc

func MakeWriter(f WriterFunc) Writer {
	return writer(f)
}

func (w writer) WriteObjekte(o Objekte) (err error) {
	return WriterFunc(w)(o)
}
