package proto_objekte

type innerSetGeneric[T any] struct {
	SetGeneric[T]
}

type MutableSetGeneric[T any] struct {
	innerSetGeneric[T]
}

func MakeMutableSetGeneric[T any](kf KeyFunc[T], es ...T) (s MutableSetGeneric[T]) {
	s.innerSetGeneric.SetGeneric = MakeSetGeneric[T](kf, es...)

	return
}

func (es MutableSetGeneric[T]) WriterAdder() WriterFunc[T] {
	return func(e T) (err error) {
		k := es.Key(e)

		if k == "" {
			err = ErrEmptyKey[T]{Element: e}
			return
		}

		es.innerSetGeneric.SetGeneric.inner[k] = e

		return
	}
}

func (es MutableSetGeneric[T]) WriterRemover() WriterFunc[T] {
	return func(e T) (err error) {
		k := es.Key(e)

		if k == "" {
			err = ErrEmptyKey[T]{Element: e}
			return
		}

		delete(es.innerSetGeneric.SetGeneric.inner, k)

		return
	}
}

func (es MutableSetGeneric[T]) Remove(es1 ...T) {
}

// func (es MutableSetGeneric[T]) RemovePrefixes(needle T) {
// 	for haystack, _ := range es.inner {
// 		if strings.HasPrefix(haystack, needle.String()) {
// 			delete(es.inner, haystack)
// 		}
// 	}
// }

func (a MutableSetGeneric[T]) Reset(b SetLike[T]) {
	for k, _ := range a.innerSetGeneric.SetGeneric.inner {
		delete(a.inner, k)
	}

	b.Each(a.WriterAdder())
}
