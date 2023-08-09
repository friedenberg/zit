package collections_value

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

type MutableSet[
	T schnittstellen.Element,
] struct {
	K schnittstellen.StringKeyer[T]
	E map[string]T
}

func (s MutableSet[T]) Len() int {
	if s.E == nil {
		return 0
	}

	return len(s.E)
}

func (a MutableSet[T]) EqualsSetLike(
	b schnittstellen.SetLike[T],
) bool {
	if b == nil {
		return false
	}

	if a.Len() != b.Len() {
		return false
	}

	for k, va := range a.E {
		vb, ok := b.Get(k)

		if !ok || !va.EqualsAny(vb) {
			return false
		}
	}

	return true
}

func (s MutableSet[T]) Key(e T) string {
	return s.K.GetKey(e)
}

func (s MutableSet[T]) Get(k string) (e T, ok bool) {
	e, ok = s.E[k]

	return
}

func (s MutableSet[T]) ContainsKey(k string) (ok bool) {
	if k == "" {
		return
	}

	_, ok = s.E[k]

	return
}

func (s MutableSet[T]) Contains(e T) (ok bool) {
	return s.ContainsKey(s.Key(e))
}

func (s MutableSet[T]) Any() (v T) {
	for _, v1 := range s.E {
		v = v1
		break
	}

	return
}

func (s MutableSet[T]) Del(v T) (err error) {
	return s.DelKey(s.Key(v))
}

func (s MutableSet[T]) DelKey(k string) (err error) {
	delete(s.E, k)
	return
}

func (s MutableSet[T]) Add(v T) (err error) {
	s.E[s.Key(v)] = v
	return
}

func (s MutableSet[T]) Elements() (out []T) {
	out = make([]T, 0, s.Len())

	for _, v := range s.E {
		out = append(out, v)
	}

	return
}

func (s MutableSet[T]) EachKey(
	wf schnittstellen.FuncIterKey,
) (err error) {
	for v := range s.E {
		if err = wf(v); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s MutableSet[T]) Each(
	wf schnittstellen.FuncIter[T],
) (err error) {
	for _, v := range s.E {
		if err = wf(v); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a MutableSet[T]) Reset() {
	for k := range a.E {
		delete(a.E, k)
	}
}

func (a MutableSet[T]) CloneSetLike() schnittstellen.SetLike[T] {
	b := MakeSet[T](a.K)

	for k, v := range a.E {
		b.E[k] = v
	}

	return b
}

func (a MutableSet[T]) CloneMutableSetLike() schnittstellen.MutableSetLike[T] {
	c := MakeMutableSet[T](a.K)
	a.Each(c.Add)
	return c
}
