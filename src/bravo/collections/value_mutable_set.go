package collections

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type mutableSetAlias[T any] struct {
	MutableSet[T]
}

type MutableValueSet[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}] struct {
	mutableSetAlias[T]
}

func MakeMutableValueSet[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](es ...T) (s MutableValueSet[T, T1]) {
	s.mutableSetAlias = mutableSetAlias[T]{
		MutableSet: MakeMutableSetGeneric(
			func(e T) string {
				return e.String()
			},
			es...,
		),
	}

	return
}

func MakeMutableValueSetStrings[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](vs ...string) (s MutableValueSet[T, T1], err error) {
	es := make([]T, len(vs))

	for i, v := range vs {
		e1 := T1(new(T))

		if err = e1.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		es[i] = T(*e1)
	}

	s = MakeMutableValueSet[T, T1](es...)

	return
}

func (es MutableValueSet[T, T1]) AddString(v string) (err error) {
	e := T1(new(T))

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	err = es.Add(*e)

	return
}

func (es MutableValueSet[T, T1]) RemovePrefixes(needle T) {
	es.Chain(
		func(e T) (err error) {
			if !strings.HasPrefix(e.String(), needle.String()) {
				err = io.EOF
			}

			return
		},
		es.Del,
	)
}

func (es MutableValueSet[T, T1]) Copy() (out ValueSet[T, T1]) {
	out.setAlias = setAlias[T]{
		Set: MakeSet[T](es.Key, es.Elements()...),
	}

	return
}
