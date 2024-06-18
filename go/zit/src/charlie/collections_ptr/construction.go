package collections_ptr

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
)

func MakeValueSetString[
	T schnittstellen.Stringer,
	TPtr schnittstellen.StringSetterPtr[T],
](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...string,
) (s Set[T, TPtr], err error) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyerPtr[T, TPtr]{}
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

func MakeValueSetValue[
	T schnittstellen.Stringer,
	TPtr schnittstellen.StringerPtr[T],
](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...T,
) (s Set[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyerPtr[T, TPtr]{}
	}

	s.K = keyer

	for i := range es {
		e := TPtr(&es[i])
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeValueSet[
	T schnittstellen.Stringer,
	TPtr schnittstellen.StringerPtr[T],
](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...TPtr,
) (s Set[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyerPtr[T, TPtr]{}
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeSetValue[T any, TPtr schnittstellen.Ptr[T]](
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

func MakeSet[T any, TPtr schnittstellen.Ptr[T]](
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

func MakeMutableValueSetValue[
	T schnittstellen.Stringer,
	TPtr schnittstellen.StringerPtr[T],
](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...T,
) (s MutableSet[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyerPtr[T, TPtr]{}
	}

	s.K = keyer

	for i := range es {
		e := TPtr(&es[i])
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeMutableValueSet[
	T schnittstellen.Stringer,
	TPtr schnittstellen.StringerPtr[T],
](
	keyer schnittstellen.StringKeyerPtr[T, TPtr],
	es ...TPtr,
) (s MutableSet[T, TPtr]) {
	s.E = make(map[string]TPtr, len(es))

	if keyer == nil {
		keyer = iter.StringerKeyerPtr[T, TPtr]{}
	}

	s.K = keyer

	for i := range es {
		e := es[i]
		s.E[s.K.GetKeyPtr(e)] = e
	}

	return
}

func MakeMutableSetValue[
	T schnittstellen.Stringer,
	TPtr schnittstellen.StringerPtr[T],
](
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

func MakeMutableSet[T any, TPtr schnittstellen.Ptr[T]](
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
