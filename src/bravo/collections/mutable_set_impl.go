package collections

import "sync"

type mutableSet[T any] struct {
	set[T]
	lock sync.Locker
}

func makeMutableSetGeneric[T any](kf KeyFunc[T], es ...T) (ms mutableSet[T]) {
	ms = mutableSet[T]{
		set:  makeSetGeneric(kf, es...),
		lock: &sync.Mutex{},
	}

	ms.set.open()

	return
}

func (es mutableSet[T]) Add(e T) (err error) {
	k := es.Key(e)

	if k == "" {
		err = ErrEmptyKey[T]{Element: e}
		return
	}

	es.lock.Lock()
	defer es.lock.Unlock()

	es.add(e)

	return
}

func (es mutableSet[T]) DelKey(k string) (err error) {
	if k == "" {
		err = ErrEmptyKey[T]{}
		return
	}

	es.lock.Lock()
	defer es.lock.Unlock()

	delete(es.set.inner, k)

	return
}

func (es mutableSet[T]) Del(e T) (err error) {
	if err = es.DelKey(es.Key(e)); err != nil {
		err = ErrEmptyKey[T]{Element: e}
		return
	}

	return
}

func (a mutableSet[T]) Reset(b SetLike[T]) {
	a.Each(a.Del)

	if b != nil {
		b.Each(a.Add)
	}
}
