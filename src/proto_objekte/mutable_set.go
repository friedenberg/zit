package proto_objekte

import (
	"strings"
)

type set[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}] struct {
	Set[T, T1]
}

type MutableSet[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}] struct {
	set[T, T1]
}

func MakeMutableSet[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](es ...T) (s MutableSet[T, T1]) {
	s.set.Set = MakeSet[T, T1](es...)

	return
}

func (es MutableSet[T, T1]) Add(e T) {
	es.addOnlyExact(e)
}

func (es MutableSet[T, T1]) addOnlyExact(e T) {
	es.inner[e.String()] = e
}

func (es MutableSet[T, T1]) Remove(es1 ...T) {
	for _, e := range es1 {
		delete(es.inner, e.String())
	}
}

func (es MutableSet[T, T1]) RemovePrefixes(needle T) {
	for haystack, _ := range es.inner {
		if strings.HasPrefix(haystack, needle.String()) {
			delete(es.inner, haystack)
		}
	}
}

func (a MutableSet[T, T1]) Equals(b MutableSet[T, T1]) bool {
	return a.Equals(b)
}

func (s1 MutableSet[T, T1]) Merge(s2 Set[T, T1]) {
	for _, e := range s2.inner {
		s1.addOnlyExact(e)
	}
}

func (s1 MutableSet[T, T1]) Copy() (s2 Set[T, T1]) {
	s2 = MakeSet[T, T1]()
	s2.open()
	defer s2.close()

	for _, e := range s1.inner {
		s2.add(e)
	}

	return
}
