package collections

func MakeWriterDoNotRepool[T any]() WriterFunc[*T] {
	return func(e *T) (err error) {
		err = ErrDoNotRepool{}
		return
	}
}

func MakeWriterNoop[T any]() WriterFunc[T] {
	return func(e T) (err error) {
		return
	}
}
