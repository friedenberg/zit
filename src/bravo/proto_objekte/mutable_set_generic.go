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

func (es MutableSetGeneric[T]) Add(e T) {
	k := es.Key(e)

	if k == "" {
		panic("empty key")
	}

	es.innerSetGeneric.SetGeneric.inner[k] = e
}

func (es MutableSetGeneric[T]) Remove(es1 ...T) {
	for _, e := range es1 {
		k := es.Key(e)

		if k == "" {
			panic("empty key")
		}

		delete(es.innerSetGeneric.SetGeneric.inner, k)
	}
}

// func (es MutableSetGeneric[T]) RemovePrefixes(needle T) {
// 	for haystack, _ := range es.inner {
// 		if strings.HasPrefix(haystack, needle.String()) {
// 			delete(es.inner, haystack)
// 		}
// 	}
// }

func (a MutableSetGeneric[T]) AddFrom(b SetLike[T]) {
	b.Each(
		func(e T) error {
			a.Add(e)
			return nil
		},
	)
}

func (a MutableSetGeneric[T]) Reset(b SetLike[T]) {
	for k, _ := range a.innerSetGeneric.SetGeneric.inner {
		delete(a.inner, k)
	}

	a.AddFrom(b)
}
