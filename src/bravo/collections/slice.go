package collections

import (
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Slice[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}] []T

func MakeSlice[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](es ...T) (s Slice[T, T1]) {
	s = make([]T, len(es))

	for i, e := range es {
		s[i] = e
	}

	return
}

func MakeSliceStrings[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](es ...string) (s Slice[T, T1], err error) {
	s = make([]T, len(es))

	for i, e := range es {
		e1 := T1(new(T))

		if err = e1.Set(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		s[i] = *e1
	}

	return
}

func (s Slice[T, T1]) Len() int {
	return len(s)
}

func (es *Slice[T, T1]) AddString(v string) (err error) {
	e1 := T1(new(T))

	if err = e1.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	es.Add(*e1)

	return
}

func (es *Slice[T, T1]) Add(e T) {
	*es = append(*es, e)
}

func (s *Slice[T, T1]) ValueSet(v string) (err error) {
	es := strings.Split(v, ",")

	for _, e := range es {
		if err = s.AddString(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (es Slice[T, T1]) SortedString() (out []string) {
	out = make([]string, len(es))

	i := 0

	for _, e := range es {
		out[i] = e.String()
		i++
	}

	sort.Slice(
		out,
		func(i, j int) bool {
			return out[i] < out[j]
		},
	)

	return
}

func (s Slice[T, T1]) String() string {
	return strings.Join(s.SortedString(), ", ")
}

func (s Slice[T, T1]) ToSet() (se ValueSet[T, T1]) {
	se = MakeSet[T, T1]()
	se.open()
	defer se.close()

	for _, e := range s {
		se.add(e)
	}

	return
}
