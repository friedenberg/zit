package collections_ptr

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
)

type Set[
	T any,
	TPtr interfaces.Ptr[T],
] struct {
	K interfaces.StringKeyerPtr[T, TPtr]
	E map[string]TPtr
}

func (s Set[T, TPtr]) Len() int {
	if s.E == nil {
		return 0
	}

	return len(s.E)
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

func (s Set[T, TPtr]) EachKey(wf interfaces.FuncIterKey) (err error) {
	for v := range s.E {
		if err = wf(v); err != nil {
			if errors.Is(err, quiter.MakeErrStopIteration()) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
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

func (s Set[T, TPtr]) Each(wf interfaces.FuncIter[T]) (err error) {
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

func (s Set[T, TPtr]) EachPtr(
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

func (a Set[T, TPtr]) CloneSetLike() interfaces.SetLike[T] {
	return a
}

func (a Set[T, TPtr]) CloneMutableSetLike() interfaces.MutableSetLike[T] {
	c := MakeMutableSet[T, TPtr](a.K)
	a.Each(c.Add)
	return c
}

func (a Set[T, TPtr]) CloneSetPtrLike() interfaces.SetPtrLike[T, TPtr] {
	return a
}

func (a Set[T, TPtr]) CloneMutableSetPtrLike() interfaces.MutableSetPtrLike[T, TPtr] {
	c := MakeMutableSet[T, TPtr](a.K)
	a.Each(c.Add)
	return c
}
