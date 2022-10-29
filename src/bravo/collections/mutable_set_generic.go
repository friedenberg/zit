package collections

import (
	"sync"
)

type innerSetGeneric[T any] struct {
	SetGeneric[T]
}

type MutableSetGeneric[T any] struct {
	innerSetGeneric[T]
	lock sync.Locker
}

func MakeMutableSetGeneric[T any](kf KeyFunc[T], es ...T) (s MutableSetGeneric[T]) {
	s.innerSetGeneric.SetGeneric = MakeSetGeneric[T](kf, es...)
	s.lock = &sync.Mutex{}

	return
}

func (es MutableSetGeneric[T]) WriterAdder() WriterFunc[T] {
	return func(e T) (err error) {
		k := es.Key(e)

		if k == "" {
			err = ErrEmptyKey[T]{Element: e}
			return
		}

		es.lock.Lock()
		defer es.lock.Unlock()
		es.innerSetGeneric.SetGeneric.inner[k] = e

		return
	}
}

func (es MutableSetGeneric[T]) WriterRemoverKeys() WriterFuncKey {
	return func(k string) (err error) {
		if k == "" {
			err = ErrEmptyKey[T]{}
			return
		}

		es.lock.Lock()
		defer es.lock.Unlock()

		delete(es.innerSetGeneric.SetGeneric.inner, k)

		return
	}
}

func (es MutableSetGeneric[T]) WriterRemover() WriterFunc[T] {
	removerKeys := es.WriterRemoverKeys()

	return func(e T) (err error) {
		if err = removerKeys(es.KeyFunc()(e)); err != nil {
			err = ErrEmptyKey[T]{Element: e}
			return
		}

		return
	}
}

// func (es MutableSetGeneric[T]) RemovePrefixes(needle T) {
// 	for haystack, _ := range es.inner {
// 		if strings.HasPrefix(haystack, needle.String()) {
// 			delete(es.inner, haystack)
// 		}
// 	}
// }

func (a MutableSetGeneric[T]) Reset(b SetLike[T]) {
	a.Each(a.WriterRemover())
	b.Each(a.WriterAdder())
}
