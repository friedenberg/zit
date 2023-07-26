package collections2

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

func MakeValueSetValue[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...T,
) (s Set[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T, TPtr]{}
	}

	s.K = keyer

	for i := range es {
		e := TPtr(&es[i])
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeValueSet[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...TPtr,
) (s Set[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T, TPtr]{}
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeSetValue[T schnittstellen.Element, TPtr schnittstellen.ElementPtr[T]](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...T,
) (s Set[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		panic("keyer was nil")
	}

	s.K = keyer

	for i := range es {
		e := TPtr(&es[i])
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeSet[T schnittstellen.Element, TPtr schnittstellen.ElementPtr[T]](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...TPtr,
) (s Set[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		panic("keyer was nil")
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeMutableValueSetValue[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...T,
) (s MutableSet[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T, TPtr]{}
	}

	s.K = keyer

	for i := range es {
		e := TPtr(&es[i])
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeMutableValueSet[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...TPtr,
) (s MutableSet[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyer[T, TPtr]{}
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeMutableSetValue[T schnittstellen.Element, TPtr schnittstellen.ElementPtr[T]](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...T,
) (s MutableSet[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		panic("keyer was nil")
	}

	s.K = keyer

	for i := range es {
		e := TPtr(&es[i])
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeMutableSet[T schnittstellen.Element, TPtr schnittstellen.ElementPtr[T]](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...TPtr,
) (s MutableSet[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		panic("keyer was nil")
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}
