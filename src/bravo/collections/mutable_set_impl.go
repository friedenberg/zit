package collections

import "sync"

type mutableSetGeneric[T any] struct {
	setGeneric[T]
	lock sync.Locker
}

func makeMutableSetGeneric[T any](kf KeyFunc[T], es ...T) (ms mutableSetGeneric[T]) {
	ms = mutableSetGeneric[T]{
		setGeneric: makeSetGeneric(kf, es...),
		lock:       &sync.Mutex{},
	}

	ms.setGeneric.open()

	return
}

func (es mutableSetGeneric[T]) Add(e T) (err error) {
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

func (es mutableSetGeneric[T]) tryLock() (ok bool) {
	if es.lock == nil {
		return
	}

	ok = true

	es.lock.Lock()

	return
}

func (es mutableSetGeneric[T]) Len() (l int) {
	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()

	l = es.setGeneric.Len()
	return
}

func (es mutableSetGeneric[T]) Get(k string) (e T, ok bool) {
	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()
	e, ok = es.setGeneric.Get(k)
	return
}

func (es mutableSetGeneric[T]) ContainsKey(k string) (ok bool) {
	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()
	ok = es.setGeneric.ContainsKey(k)
	return
}

func (es mutableSetGeneric[T]) Contains(e T) (ok bool) {
	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()
	ok = es.setGeneric.Contains(e)
	return
}

func (es mutableSetGeneric[T]) DelKey(k string) (err error) {
	if k == "" {
		err = ErrEmptyKey[T]{}
		return
	}

	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()

	delete(es.setGeneric.inner, k)

	return
}

func (es mutableSetGeneric[T]) Del(e T) (err error) {
	if err = es.DelKey(es.Key(e)); err != nil {
		err = ErrEmptyKey[T]{Element: e}
		return
	}

	return
}

func (es mutableSetGeneric[T]) Each(wf WriterFunc[T]) (err error) {
	if !es.tryLock() {
		return
	}

	defer es.lock.Unlock()

	return es.setGeneric.Each(wf)
}

func (a mutableSetGeneric[T]) Reset(b SetLike[T]) {
	a.Each(a.Del)
	b.Each(a.Add)
}
