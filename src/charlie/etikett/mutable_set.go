package etikett

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type innerSet = Set

type MutableSet struct {
	innerSet
}

func MakeMutableSet(es ...Etikett) (s MutableSet) {
	s.innerSet = innerSet(MakeSet(es...))

	return
}

func (es MutableSet) AddString(v string) (err error) {
	var e Etikett

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err) 
    return
	}

	es.addOnlyExact(e)

	return
}

func (es MutableSet) AddNormalized(e Etikett) {
	expanded := e.Expanded(ExpanderRight{})
	i := MakeSet(e)
  i.open()
  defer i.close()

	for _, e := range expanded.inner {
		if es.Contains(e) {
			es.Remove(e)
			i.addOnlyExact(e)
		}
	}

	for _, e1 := range i.WithRemovedCommonPrefixes().inner {
		es.addOnlyExact(e1)
	}
}

func (es MutableSet) Add(e Etikett) {
	// expanded := e.Expanded(ExpanderRight{})
	// intersection := es.Intersect(expanded)
	// es.Remove(intersection.Etiketten()...)

	es.addOnlyExact(e)
}

func (es MutableSet) addOnlyExact(e Etikett) {
	es.inner[e.String()] = e
}

func (s MutableSet) Set(v string) (err error) {
	es := strings.Split(v, ",")

	for _, e := range es {
		if err = s.AddString(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (es MutableSet) Remove(es1 ...Etikett) {
	for _, e := range es1 {
		delete(es.inner, e.String())
	}
}

func (es MutableSet) RemovePrefixes(needle Etikett) {
	for haystack, _ := range es.inner {
		if strings.HasPrefix(haystack, needle.String()) {
			delete(es.inner, haystack)
		}
	}
}

func (a MutableSet) Equals(b MutableSet) bool {
	return a.innerSet.Equals(b.innerSet)
}

func (s1 MutableSet) Merge(s2 Set) {
	for _, e := range s2.inner {
		s1.addOnlyExact(e)
	}
}

func (s1 MutableSet) Copy() (s2 Set) {
	s2 = MakeSet()
  s2.open()
  defer s2.close()

	for _, e := range s1.inner {
		s2.addOnlyExact(e)
	}

	return
}
