package collections_ptr

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

type Set[
	T schnittstellen.Element,
	TPtr schnittstellen.ElementPtr[T],
] struct {
	K schnittstellen.StringKeyerPtr[T, TPtr]
	E map[string]TPtr
}

func (s Set[T, TPtr]) Len() int {
	if s.E == nil {
		return 0
	}

	return len(s.E)
}

func (a Set[T, TPtr]) EqualsSetPtrLike(
	b schnittstellen.SetPtrLike[T, TPtr],
) bool {
	return a.EqualsSetLike(b)
}

func (a Set[T, TPtr]) EqualsSetLike(b schnittstellen.SetLike[T]) bool {
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

func (s Set[T, TPtr]) Key(e T) string {
	return s.K.GetKey(e)
}

func (s Set[T, TPtr]) KeyPtr(e TPtr) string {
	return s.K.GetKeyPtr(e)
}

func (s Set[T, TPtr]) GetPtr(k string) (e TPtr, ok bool) {
	e, ok = s.E[k]

	return
}

func (s Set[T, TPtr]) Get(k string) (e T, ok bool) {
	var e1 TPtr

	if e1, ok = s.E[k]; ok {
		e = *e1
	}

	return
}

func (s Set[T, TPtr]) Any() (e T) {
	for _, e1 := range s.E {
		return *e1
	}

	return
}

func (s Set[T, TPtr]) ContainsKey(k string) (ok bool) {
	if k == "" {
		return
	}

	_, ok = s.E[k]

	return
}

func (s Set[T, TPtr]) Contains(e T) (ok bool) {
	return s.ContainsKey(s.Key(e))
}

func (s Set[T, TPtr]) EachKey(wf schnittstellen.FuncIterKey) (err error) {
	for v := range s.E {
		if err = wf(v); err != nil {
			if errors.Is(err, iter.MakeErrStopIteration()) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s Set[T, TPtr]) Elements() (out []T) {
	out = make([]T, 0, s.Len())

	for _, v := range s.E {
		out = append(out, *v)
	}

	return
}

func (s Set[T, TPtr]) Add(v T) (err error) {
	s.E[s.Key(v)] = TPtr(&v)
	return
}

func (s Set[T, TPtr]) AddPtr(v TPtr) (err error) {
	s.E[s.K.GetKeyPtr(v)] = v
	return
}

func (s Set[T, TPtr]) Each(wf schnittstellen.FuncIter[T]) (err error) {
	for _, v := range s.E {
		if err = wf(*v); err != nil {
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

func (s Set[T, TPtr]) EachPtr(
	wf schnittstellen.FuncIter[TPtr],
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

func (a Set[T, TPtr]) CloneSetLike() schnittstellen.SetLike[T] {
	return a
}

func (a Set[T, TPtr]) CloneMutableSetLike() schnittstellen.MutableSetLike[T] {
	c := MakeMutableSet[T, TPtr](a.K)
	a.Each(c.Add)
	return c
}

func (a Set[T, TPtr]) CloneSetPtrLike() schnittstellen.SetPtrLike[T, TPtr] {
	return a
}

func (a Set[T, TPtr]) CloneMutableSetPtrLike() schnittstellen.MutableSetPtrLike[T, TPtr] {
	c := MakeMutableSet[T, TPtr](a.K)
	a.Each(c.Add)
	return c
}