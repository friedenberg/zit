package collections

import (
	"fmt"
	"reflect"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type set[T any] struct {
	keyFunc func(T) string
	closed  bool
	inner   map[string]T
}

type setAlias[T any] struct {
	Set[T]
}

func makeSet[T any](kf KeyFunc[T], es ...T) (s set[T]) {
	t := *new(T)
	// Required because interface types do not properly get handled by
	// `reflect.TypeOf`
	t1 := make([]T, 1)

	if reflect.TypeOf(t1).Elem().Kind() == reflect.Interface {
		kf(t1[0])
	} else {
		// confirms that the key function supports nil pointers properly
		switch reflect.TypeOf(t).Kind() {
		// case reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		case reflect.Ptr:
			kf(t)
		}
	}

	s.keyFunc = kf
	s.inner = make(map[string]T, len(es))
	s.open()
	defer s.close()

	for _, e := range es {
		s.add(e)
	}

	return
}

func (s *set[T]) open() {
	s.closed = false
}

func (s *set[T]) close() {
	s.closed = true
}

func (s set[T]) Len() int {
	if s.inner == nil {
		return 0
	}

	return len(s.inner)
}

func (s set[T]) Key(e T) string {
	return s.keyFunc(e)
}

func (s set[T]) Get(k string) (e T, ok bool) {
	e, ok = s.inner[k]
	return
}

func (s set[T]) ContainsKey(k string) (ok bool) {
	if k == "" {
		return
	}

	_, ok = s.inner[k]

	return
}

func (s set[T]) Contains(e T) (ok bool) {
	return s.ContainsKey(s.Key(e))
}

func (es set[T]) add(e T) (err error) {
	if es.closed {
		panic(fmt.Sprintf("trying to add %T to closed set", e))
	}

	es.inner[es.Key(e)] = e

	return
}

func (s set[T]) EachKey(wf WriterFuncKey) (err error) {
	for v := range s.inner {
		if err = wf(v); err != nil {
			if errors.Is(err, MakeErrStopIteration()) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s set[T]) Each(wf WriterFunc[T]) (err error) {
	for _, v := range s.inner {
		if err = wf(v); err != nil {
			if errors.Is(err, MakeErrStopIteration()) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s set[T]) EachPtr(wf WriterFunc[*T]) (err error) {
	for _, v := range s.inner {
		if err = wf(&v); err != nil {
			if errors.Is(err, MakeErrStopIteration()) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}
