package collections

func MakeWriterNil[T any]() WriterFunc[*T] {
	return func(e *T) (err error) {
		e = nil
		return
	}
}

func MakeWriterNoop[T any]() WriterFunc[T] {
	return func(e T) (err error) {
		return
	}
}
