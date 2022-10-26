package objekte

type WriterWithIndex interface {
	WriteObjekteWithIndex(ObjekteWithIndex) error
}

type WriterWithIndexFunc func(ObjekteWithIndex) error

type writerWithIndex WriterWithIndexFunc

func MakeWriterWithIndex(f WriterWithIndexFunc) WriterWithIndex {
	return writerWithIndex(f)
}

func (w writerWithIndex) WriteObjekteWithIndex(o ObjekteWithIndex) (err error) {
	return WriterWithIndexFunc(w)(o)
}
