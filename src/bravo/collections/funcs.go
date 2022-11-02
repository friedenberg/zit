package collections

func MakeNiller[T any]() WriterFunc[*T] {
	return func(e *T) (err error) {
		e = nil
		return
	}
}
