package collections_value

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

func MakeValueSetString[
	T schnittstellen.ValueLike,
	TPtr interface {
		schnittstellen.ValuePtr[T]
		schnittstellen.Setter
	},
](
	keyer schnittstellen.StringKeyer[T],
	es ...string,
) (s Set[T], err error) {
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}.RegisterGob()
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

func MakeValueSetValue[T schnittstellen.ValueLike](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s Set[T]) {
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}.RegisterGob()
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeValueSet[T schnittstellen.ValueLike](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s Set[T]) {
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}.RegisterGob()
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeSetValue[T schnittstellen.Element](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s Set[T]) {
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

func MakeSet[T schnittstellen.Element](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s Set[T]) {
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

func MakeMutableValueSetValue[T schnittstellen.ValueLike](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s MutableSet[T]) {
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}.RegisterGob()
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeMutableValueSet[T schnittstellen.ValueLike](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s MutableSet[T]) {
	s.E = make(map[string]T, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T]{}.RegisterGob()
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKey(e)] = e
	}

	return
}

func MakeMutableSetValue[T schnittstellen.Element](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s MutableSet[T]) {
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

func MakeMutableSet[T schnittstellen.Element](
	keyer schnittstellen.StringKeyer[T],
	es ...T,
) (s MutableSet[T]) {
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
