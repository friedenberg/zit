package collections

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type setAlias[T any] struct {
	Set[T]
}

type ValueSet[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}] struct {
	setAlias[T]
}

func MakeValueSet[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](es ...T) (s ValueSet[T, T1]) {
	s.setAlias = setAlias[T]{
		Set: MakeSetGeneric(
			func(e T) string {
				return e.String()
			},
			es...,
		),
	}

	return
}

func MakeValueSetStrings[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](vs ...string) (s ValueSet[T, T1], err error) {
	es := make([]T, len(vs))

	for i, v := range vs {
		e1 := T1(new(T))

		if err = e1.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		es[i] = T(*e1)
	}

	s = MakeValueSet[T, T1](es...)

	return
}

func (s *ValueSet[T, T1]) Set(v string) (err error) {
	parts := strings.Split(v, ",")

	if *s, err = MakeValueSetStrings[T, T1](parts...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s1 ValueSet[T, T1]) Copy() (s2 ValueSet[T, T1]) {
	s2 = MakeValueSet[T, T1](s1.Elements()...)

	return
}

func (s1 ValueSet[T, T1]) MutableCopy() (s2 MutableValueSet[T, T1]) {
	s2 = MakeMutableValueSet[T, T1]()
	s1.Each(s2.Add)

	return
}

func (s ValueSet[T, T1]) Strings() (out []string) {
	out = make([]string, 0, s.Len())

	s.Each(
		func(e T) (err error) {
			out = append(out, e.String())

			return
		},
	)

	return
}

func (es ValueSet[T, T1]) Sorted() (out []T) {
	out = es.Elements()

	sort.Slice(
		out,
		func(i, j int) bool {
			return out[i].String() < out[j].String()
		},
	)

	return
}

func (es ValueSet[T, T1]) SortedString() (out []string) {
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

func (s ValueSet[T, T1]) ContainsString(es string) bool {
	return s.ContainsKey(es)
}

func (s1 ValueSet[T, T1]) IntersectPrefixes(s2 ValueSet[T, T1]) (out ValueSet[T, T1]) {
	s3 := MakeMutableValueSet[T, T1]()

	s2.Each(
		func(e1 T) (err error) {
			s1.Each(
				func(e T) (err error) {
					if strings.HasPrefix(e.String(), e1.String()) {
						s3.Add(e)
					}

					return
				},
			)

			return
		},
	)

	out = s3.Copy()

	return
}

func (es ValueSet[T, T1]) Description() string {
	sb := &strings.Builder{}
	first := true

	for _, e1 := range es.SortedString() {
		if !first {
			sb.WriteString(", ")
		}

		sb.WriteString(e1)

		first = false
	}

	return sb.String()
}

func (s ValueSet[T, T1]) String() string {
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

func (es ValueSet[T, T1]) MarshalJSON() ([]byte, error) {
	return json.Marshal(es.SortedString())
}

func (es *ValueSet[T, T1]) UnmarshalJSON(b []byte) (err error) {
	var vs []string

	if err = json.Unmarshal(b, &vs); err != nil {
		err = errors.Wrap(err)
		return
	}

	if *es, err = MakeValueSetStrings[T, T1](vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s ValueSet[T, T1]) MarshalBinary() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *ValueSet[T, T1]) UnmarshalBinary(text []byte) (err error) {
	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
