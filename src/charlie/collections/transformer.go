package collections

func MakeFuncTransformer[T any, T1 any](wf WriterFunc[T]) WriterFunc[T1] {
	return func(e T1) (err error) {
		if e1, ok := any(e).(T); ok {
			return wf(e1)
		}

		return
	}
}
