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

func MakeMutableSetGeneric[T any](kf KeyFunc[T], es ...T) (s MutableSetLike[T]) {
	var s1 MutableSetGeneric[T]
	s1.innerSetGeneric.SetGeneric = MakeSetGeneric[T](kf, es...)
	s1.lock = &sync.Mutex{}

	s = s1

	return
}

func (es MutableSetGeneric[T]) Add(e T) (err error) {
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

func (es MutableSetGeneric[T]) tryLock() (ok bool) {
	if es.lock == nil {
		return
	}

	ok = true

	es.lock.Lock()

	return
}

func (es MutableSetGeneric[T]) Len() (l int) {
	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()

	l = es.innerSetGeneric.SetGeneric.Len()
	return
}

func (es MutableSetGeneric[T]) Get(k string) (e T, ok bool) {
	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()
	e, ok = es.innerSetGeneric.SetGeneric.Get(k)
	return
}

func (es MutableSetGeneric[T]) ContainsKey(k string) (ok bool) {
	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()
	ok = es.innerSetGeneric.SetGeneric.ContainsKey(k)
	return
}

func (es MutableSetGeneric[T]) Contains(e T) (ok bool) {
	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()
	ok = es.innerSetGeneric.SetGeneric.Contains(e)
	return
}

func (es MutableSetGeneric[T]) DelKey(k string) (err error) {
	if k == "" {
		err = ErrEmptyKey[T]{}
		return
	}

	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()

	delete(es.innerSetGeneric.SetGeneric.inner, k)

	return
}

func (es MutableSetGeneric[T]) Del(e T) (err error) {
	if err = es.DelKey(es.Key(e)); err != nil {
		err = ErrEmptyKey[T]{Element: e}
		return
	}

	return
}

func (es MutableSetGeneric[T]) Each(wf WriterFunc[T]) (err error) {
	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()

	return es.innerSetGeneric.SetGeneric.Each(wf)
}

func (a MutableSetGeneric[T]) Reset(b SetLike[T]) {
	a.Each(a.Del)
	b.Each(a.Add)
}
