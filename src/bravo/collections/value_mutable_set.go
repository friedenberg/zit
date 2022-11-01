package collections

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type set[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}] struct {
	ValueSet[T, T1]
}

type ValueMutableSet[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}] struct {
	set[T, T1]
}

func MakeMutableSet[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](es ...T) (s ValueMutableSet[T, T1]) {
	s.set.ValueSet = MakeSet[T, T1](es...)

	return
}

func MakeMutableSetStrings[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](es ...string) (s ValueMutableSet[T, T1], err error) {
	if s.set.ValueSet, err = MakeSetStrings[T, T1](es...); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (es ValueMutableSet[T, T1]) Add(e T) {
	es.inner[e.String()] = e
}

func (es ValueMutableSet[T, T1]) AddString(v string) (err error) {
	e := T1(new(T))

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	es.Add(*e)

	return
}

func (es ValueMutableSet[T, T1]) Remove(es1 ...T) {
	for _, e := range es1 {
		delete(es.inner, e.String())
	}
}

func (es ValueMutableSet[T, T1]) RemovePrefixes(needle T) {
	for haystack, _ := range es.inner {
		if strings.HasPrefix(haystack, needle.String()) {
			delete(es.inner, haystack)
		}
	}
}

func (a ValueMutableSet[T, T1]) Equals(b ValueMutableSet[T, T1]) bool {
	return a.set.ValueSet.Equals(b.set.ValueSet)
}

func (s1 ValueMutableSet[T, T1]) Merge(s2 ValueSet[T, T1]) {
	for _, e := range s2.inner {
		s1.Add(e)
	}
}

func (s1 ValueMutableSet[T, T1]) Reset(s2 ValueSet[T, T1]) {
	for k, _ := range s1.inner {
		delete(s1.inner, k)
	}

	for k, e := range s2.inner {
		s1.inner[k] = e
	}
}

func (s1 ValueMutableSet[T, T1]) Copy() (s2 ValueSet[T, T1]) {
	s2 = MakeSet[T, T1]()
	s2.open()
	defer s2.close()

	for _, e := range s1.inner {
		s2.add(e)
	}

	return
}
