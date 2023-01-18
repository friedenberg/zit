package collections

import (
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type mutableSetAlias[T any] struct {
	MutableSet[T]
}

type MutableValueSet[
	T schnittstellen.Value,
	T1 schnittstellen.ValuePtr[T],
] struct {
	mutableSetAlias[T]
}

func MakeMutableValueSet[
	T schnittstellen.Value,
	T1 schnittstellen.ValuePtr[T],
](
	es ...T,
) (s MutableValueSet[T, T1]) {
	s.mutableSetAlias = mutableSetAlias[T]{
		MutableSet: MakeMutableSet(
			func(e T) string {
				return e.String()
			},
			es...,
		),
	}

	return
}

func MakeMutableValueSetStrings[
	T schnittstellen.Value,
	T1 schnittstellen.ValuePtr[T],
](
	vs ...string,
) (s MutableValueSet[T, T1], err error) {
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

func (s MutableValueSet[T, T1]) Strings() (out []string) {
	out = make([]string, 0, s.Len())

	s.Each(
		func(e T) (err error) {
			out = append(out, e.String())

			return
		},
	)

	return
}

func (s MutableValueSet[T, T1]) StringAdder() schnittstellen.FuncSetString {
	return MakeFuncSetString[T, T1](s)
}

func (es MutableValueSet[T, T1]) RemovePrefixes(needle T) {
	es.Chain(
		func(e T) (err error) {
			if !strings.HasPrefix(e.String(), needle.String()) {
				err = ErrStopIteration
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

func (es MutableValueSet[T, T1]) SortedString() (out []string) {
	out = make([]string, 0, es.Len())

	es.Each(
		func(e T) (err error) {
			out = append(out, e.String())

			return
		},
	)

	sort.Slice(
		out,
		func(i, j int) bool {
			return out[i] < out[j]
		},
	)

	return
}

func (s MutableValueSet[T, T1]) String() string {
	if s.SetLike == nil || s.Len() == 0 {
		return ""
	}

	sb := &strings.Builder{}
	first := true

	for _, e1 := range s.SortedString() {
		if !first {
			sb.WriteString(", ")
		}

		sb.WriteString(e1)

		first = false
	}

	return sb.String()
}

func (s *MutableValueSet[T, T1]) Set(v string) (err error) {
	parts := strings.Split(v, ",")

	if len(parts) == 1 && parts[0] == "" {
		parts = []string{}
	}

	if *s, err = MakeMutableValueSetStrings[T, T1](parts...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s MutableValueSet[T, T1]) MarshalBinary() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *MutableValueSet[T, T1]) UnmarshalBinary(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
