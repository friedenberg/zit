package collections_value

import (
	"iter"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
)

type Set[
	T any,
] struct {
	K interfaces.StringKeyer[T]
	E map[string]T
}

func (s Set[T]) AllKeys() iter.Seq[string] {
	return func(yield func(string) bool) {
		for k := range s.E {
			if !yield(k) {
				break
			}
		}
	}
}

func (s Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, e := range s.E {
			if !yield(e) {
				break
			}
		}
	}
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

func (s Set[T]) EachKey(wf interfaces.FuncIterKey) (err error) {
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

func (s Set[T]) Add(v T) (err error) {
	s.E[s.Key(v)] = v
	return
}

func (s Set[T]) Each(wf interfaces.FuncIter[T]) (err error) {
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

func (a Set[T]) CloneSetLike() interfaces.SetLike[T] {
	return a
}

func (a Set[T]) CloneMutableSetLike() interfaces.MutableSetLike[T] {
	c := MakeMutableSet[T](a.K)
	a.Each(c.Add)
	return c
}
