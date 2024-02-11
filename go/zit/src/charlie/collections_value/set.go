package collections_value

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
)

type Set[
	T any,
] struct {
	K schnittstellen.StringKeyer[T]
	E map[string]T
}

func (s Set[T]) Len() int {
	if s.E == nil {
		return 0
	}

	return len(s.E)
}

func (s Set[T]) Key(e T) string {
	return s.K.GetKey(e)
}

func (s Set[T]) Get(k string) (e T, ok bool) {
	e, ok = s.E[k]

	return
}

func (s Set[T]) Any() (e T) {
	for _, e1 := range s.E {
		return e1
	}

	return
}

func (s Set[T]) ContainsKey(k string) (ok bool) {
	if k == "" {
		return
	}

	_, ok = s.E[k]

	return
}

func (s Set[T]) Contains(e T) (ok bool) {
	return s.ContainsKey(s.Key(e))
}

func (s Set[T]) EachKey(wf schnittstellen.FuncIterKey) (err error) {
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

func (s Set[T]) Add(v T) (err error) {
	s.E[s.Key(v)] = v
	return
}

func (s Set[T]) Each(wf schnittstellen.FuncIter[T]) (err error) {
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

func (a Set[T]) CloneSetLike() schnittstellen.SetLike[T] {
	return a
}

func (a Set[T]) CloneMutableSetLike() schnittstellen.MutableSetLike[T] {
	c := MakeMutableSet[T](a.K)
	a.Each(c.Add)
	return c
}
