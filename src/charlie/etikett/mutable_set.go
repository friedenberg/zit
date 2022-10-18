package etikett

import (
	"strings"

	"github.com/friedenberg/zit/src/bravo/proto_objekte"
)

type MutableSet = proto_objekte.MutableSet[Etikett, *Etikett]

func MakeMutableSet(hs ...Etikett) MutableSet {
	return MutableSet(proto_objekte.MakeMutableSet[Etikett, *Etikett](hs...))
}

func AddNormalized(es MutableSet, e Etikett) {
	es.Merge(e.Expanded(ExpanderRight{}))
	es.Add(e)

	es.Reset(WithRemovedCommonPrefixes(es.Copy()))
}

func RemovePrefixes(es MutableSet, needle Etikett) {
	for _, haystack := range es.Elements() {
		//TODO make more efficient
		if strings.HasPrefix(haystack.String(), needle.String()) {
			es.Remove(haystack)
		}
	}
}

func Withdraw(s1 MutableSet, e Etikett) (s2 Set) {
	s3 := MakeMutableSet()

	for _, e1 := range s1.Elements() {
		if e1.Contains(e) {
			s3.Add(e1)
		}
	}

	s1.Remove(s3.Elements()...)
	s2 = s3.Copy()

	return
}

// func (s MutableSet) Set(v string) (err error) {
// 	es := strings.Split(v, ",")

// 	for _, e := range es {
// 		if err = s.AddString(e); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	return
// }

// func (es MutableSet) RemovePrefixes(needle Etikett) {
// 	for haystack, _ := range es.inner {
// 		if strings.HasPrefix(haystack, needle.String()) {
// 			delete(es.inner, haystack)
// 		}
// 	}
// }
