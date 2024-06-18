package collections_value

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
)

func MakeValueSetString[
	T schnittstellen.Stringer,
	TPtr schnittstellen.StringSetterPtr[T],
](
	keyer schnittstellen.StringKeyer[T],
	es ...string,
) (s Set[T], err error) {
	gob.Register(s)
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}
	}

	s.K = keyer

	for _, v := range es {
		var e T
		e1 := TPtr(&e)

		if err = e1.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeValueSetValue[T schnittstellen.Stringer](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s Set[T]) {
	gob.Register(s)
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeValueSet[T schnittstellen.Stringer](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s Set[T]) {
	gob.Register(s)
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeSetValue[T schnittstellen.Stringer](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s Set[T]) {
	gob.Register(s)
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		panic("keyer was nil")
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeSet[T any](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s Set[T]) {
	gob.Register(s)
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		panic("keyer was nil")
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeMutableValueSetValue[T schnittstellen.Stringer](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s MutableSet[T]) {
	gob.Register(s)
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeMutableValueSet[T schnittstellen.Stringer](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s MutableSet[T]) {
	gob.Register(s)
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeMutableSetValue[T any](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s MutableSet[T]) {
	gob.Register(s)
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		panic("keyer was nil")
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeMutableSet[T any](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s MutableSet[T]) {
	gob.Register(s)
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		panic("keyer was nil")
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}
