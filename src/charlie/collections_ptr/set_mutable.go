package collections_ptr

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

type MutableSet[
	T any,
	TPtr schnittstellen.Ptr[T],
] struct {
	K schnittstellen.StringKeyerPtr[T, TPtr]
	E map[string]TPtr
}

func (s MutableSet[T, TPtr]) Len() int {
	if s.E == nil {
		return 0
	}

	return len(s.E)
}

func (s MutableSet[T, TPtr]) Key(e T) string {
	return s.K.GetKey(e)
}

func (s MutableSet[T, TPtr]) KeyPtr(e TPtr) string {
	return s.K.GetKeyPtr(e)
}

func (s MutableSet[T, TPtr]) GetPtr(k string) (e TPtr, ok bool) {
	e, ok = s.E[k]

	return
}

func (s MutableSet[T, TPtr]) Get(k string) (e T, ok bool) {
	var e1 TPtr

	if e1, ok = s.E[k]; ok {
		e = *e1
	}

	return
}

func (s MutableSet[T, TPtr]) ContainsKey(k string) (ok bool) {
	if k == "" {
		return
	}

	_, ok = s.E[k]

	return
}

func (s MutableSet[T, TPtr]) Contains(e T) (ok bool) {
	return s.ContainsKey(s.Key(e))
}

func (s MutableSet[T, TPtr]) Any() (v T) {
	for _, v1 := range s.E {
		v = *v1
		break
	}

	return
}

func (s MutableSet[T, TPtr]) Del(v T) (err error) {
	return s.DelKey(s.Key(v))
}

func (s MutableSet[T, TPtr]) DelPtr(v TPtr) (err error) {
	return s.DelKey(s.K.GetKeyPtr(v))
}

func (s MutableSet[T, TPtr]) DelKey(k string) (err error) {
	delete(s.E, k)
	return
}

func (s MutableSet[T, TPtr]) Add(v T) (err error) {
	s.E[s.Key(v)] = TPtr(&v)
	return
}

func (s MutableSet[T, TPtr]) AddPtr(v TPtr) (err error) {
	s.E[s.K.GetKeyPtr(v)] = v
	return
}

func (s MutableSet[T, TPtr]) Elements() (out []T) {
	out = make([]T, 0, s.Len())

	for _, v := range s.E {
		out = append(out, *v)
	}

	return
}

func (s MutableSet[T, TPtr]) EachKey(
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

func (s MutableSet[T, TPtr]) Each(
	wf schnittstellen.FuncIter[T],
) (err error) {
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

func (s MutableSet[T, TPtr]) EachPtr(
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

func (a MutableSet[T, TPtr]) Reset() {
	for k := range a.E {
		delete(a.E, k)
	}
}

func (a MutableSet[T, TPtr]) CloneSetLike() schnittstellen.SetLike[T] {
	b := MakeSet[T, TPtr](a.K)

	for k, v := range a.E {
		b.E[k] = v
	}

	return b
}

func (a MutableSet[T, TPtr]) CloneMutableSetLike() schnittstellen.MutableSetLike[T] {
	c := MakeMutableSet[T, TPtr](a.K)
	a.Each(c.Add)
	return c
}

func (a MutableSet[T, TPtr]) CloneSetPtrLike() schnittstellen.SetPtrLike[T, TPtr] {
	b := MakeSet[T, TPtr](a.K)

	for k, v := range a.E {
		b.E[k] = v
	}

	return b
}

func (a MutableSet[T, TPtr]) CloneMutableSetPtrLike() schnittstellen.MutableSetPtrLike[T, TPtr] {
	c := MakeMutableSet[T, TPtr](a.K)
	a.Each(c.Add)
	return c
}
