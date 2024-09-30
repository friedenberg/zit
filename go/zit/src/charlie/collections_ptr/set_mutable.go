package collections_ptr

import (
	"iter"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
)

type MutableSet[
	T any,
	TPtr interfaces.Ptr[T],
] struct {
	K interfaces.StringKeyerPtr[T, TPtr]
	E map[string]TPtr
}

func (s MutableSet[T, TPtr]) AllKeys() iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range s.E {
			if !yield(k) {
				break
			}
		}
	}
}

func (s MutableSet[T, TPtr]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, e := range s.E {
			if !yield(*e) {
				break
			}
		}
	}
}

func (s MutableSet[T, TPtr]) AllPtr() iter.Seq[TPtr] {
	return func(yield func(TPtr) bool) {
		for _, e := range s.E {
			if !yield(e) {
				break
			}
		}
	}
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

func (s MutableSet[T, TPtr]) EachKey(
	wf interfaces.FuncIterKey,
) (err error) {
	for v := range s.E {
		if err = wf(v); err != nil {
			if quiter.IsStopIteration(err) {
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
	wf interfaces.FuncIter[T],
) (err error) {
	for _, v := range s.E {
		if err = wf(*v); err != nil {
			if quiter.IsStopIteration(err) {
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
	wf interfaces.FuncIter[TPtr],
) (err error) {
	for _, v := range s.E {
		if err = wf(v); err != nil {
			if quiter.IsStopIteration(err) {
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

func (a MutableSet[T, TPtr]) CloneSetLike() interfaces.SetLike[T] {
	b := MakeSet[T, TPtr](a.K)

	for k, v := range a.E {
		b.E[k] = v
	}

	return b
}

func (a MutableSet[T, TPtr]) CloneMutableSetLike() interfaces.MutableSetLike[T] {
	c := MakeMutableSet[T, TPtr](a.K)
	a.Each(c.Add)
	return c
}

func (a MutableSet[T, TPtr]) CloneSetPtrLike() interfaces.SetPtrLike[T, TPtr] {
	b := MakeSet[T, TPtr](a.K)

	for k, v := range a.E {
		b.E[k] = v
	}

	return b
}

func (a MutableSet[T, TPtr]) CloneMutableSetPtrLike() interfaces.MutableSetPtrLike[T, TPtr] {
	c := MakeMutableSet[T, TPtr](a.K)
	a.Each(c.Add)
	return c
}
