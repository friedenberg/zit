package collections

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type mutableSet[T schnittstellen.ValueLike] struct {
	set[T]
	lock sync.Locker
}

func MakeMutableSetStringer[T schnittstellen.ValueLike](
	es ...T,
) schnittstellen.MutableSetLike[T] {
	return MakeMutableSet(
		(T).String,
		es...,
	)
}

func MakeMutableSet[T schnittstellen.ValueLike](
	kf KeyFunc[T],
	es ...T,
) schnittstellen.MutableSetLike[T] {
	s := makeSet(kf, es...)

	ms := &mutableSet[T]{
		set:  *s,
		lock: &sync.Mutex{},
	}

	ms.set.open()

	return ms
}

func (es mutableSet[T]) AddCustomKey(e T, kf func(T) string) (err error) {
	k := kf(e)

	if k == "" {
		err = errors.Wrap(ErrEmptyKey[T]{Element: e})
		return
	}

	es.lock.Lock()
	defer es.lock.Unlock()

	es.addCustom(e, kf)

	return
}

func (es mutableSet[T]) Add(e T) (err error) {
	k := es.Key(e)

	if k == "" {
		err = errors.Wrap(ErrEmptyKey[T]{Element: e})
		return
	}

	es.lock.Lock()
	defer es.lock.Unlock()

	es.add(e)

	return
}

func (es mutableSet[T]) DelKey(k string) (err error) {
	if k == "" {
		err = errors.Wrap(ErrEmptyKey[T]{})
		return
	}

	es.lock.Lock()
	defer es.lock.Unlock()

	delete(es.set.elementMap, k)

	return
}

func (es mutableSet[T]) Del(e T) (err error) {
	if err = es.DelKey(es.Key(e)); err != nil {
		err = errors.Wrap(ErrEmptyKey[T]{Element: e})
		return
	}

	return
}

func (a *mutableSet[T]) Reset() {
	a.Each(a.Del)
	a.lock = &sync.Mutex{}
}

func (a mutableSet[T]) CloneSetLike() schnittstellen.SetLike[T] {
	return a.set.CloneSetLike()
}

func (a mutableSet[T]) CloneMutableSetLike() schnittstellen.MutableSetLike[T] {
	return a.set.CloneMutableSetLike()
}

func (s mutableSet[T]) MarshalBinary() (bs []byte, err error) {
	b := bytes.NewBuffer(bs)
	enc := gob.NewEncoder(b)

	if err = enc.Encode(s.set.elementMap); err != nil {
		err = errors.Wrap(err)
		return
	}

	bs = b.Bytes()

	return
}

func (s *mutableSet[T]) UnmarshalBinary(bs []byte) (err error) {
	s.Reset()

	b := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(b)

	if err = dec.Decode(&s.set.elementMap); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
