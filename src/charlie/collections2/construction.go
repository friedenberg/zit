package collections2

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
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...string,
) (s Set[T, TPtr], err error) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyerPtr[T, TPtr]{}.RegisterGob()
	}

	s.K = keyer

	for _, v := range es {
		var e T
		e1 := TPtr(&e)

		if err = e1.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		s.E[s.K.GetKeyPtr(e1)] = e1
	}

	return
}

func MakeValueSetValue[T schnittstellen.ValueLike, TPtr schnittstellen.ValuePtr[T]](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...T,
) (s Set[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyerPtr[T, TPtr]{}.RegisterGob()
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
		keyer = iter.StringerKeyerPtr[T, TPtr]{}.RegisterGob()
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
		keyer = iter.StringerKeyerPtr[T, TPtr]{}.RegisterGob()
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
		keyer = iter.StringerKeyerPtr[T, TPtr]{}.RegisterGob()
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
