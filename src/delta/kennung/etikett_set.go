package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type EtikettSet = collections.ValueSet[Etikett, *Etikett]

func MakeEtikettSet(es ...Etikett) (s EtikettSet) {
	return EtikettSet(collections.MakeValueSet(es...))
}

func MakeSetStrings(vs ...string) (s EtikettSet, err error) {
	var s1 collections.ValueSet[Etikett, *Etikett]

	if s1, err = collections.MakeValueSetStrings[Etikett, *Etikett](vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	s = EtikettSet(s1)

	return
}

func WithRemovedCommonPrefixes(s EtikettSet) (s2 EtikettSet) {
	es1 := s.Sorted()
	es := make([]Etikett, 0, len(es1))

	for _, e := range es1 {
		if len(es) == 0 {
			es = append(es, e)
			continue
		}

		idxLast := len(es) - 1
		last := es[idxLast]

		switch {
		case last.Contains(e):
			continue

		case e.Contains(last):
			es[idxLast] = e

		default:
			es = append(es, e)
		}
	}

	s2 = MakeEtikettSet(es...)

	return
}

func Expanded(s EtikettSet, exes ...Expander) (out EtikettSet) {
	s1 := MakeEtikettMutableSet()

	for _, e := range s.Elements() {
		e.Expanded(exes...).Each(s1.Add)
	}

	out = s1.Copy()

	return
}

// func (s Set) String() string {
// 	return strings.Join(s.SortedString(), ", ")
// }

func IntersectPrefixes(s1 EtikettSet, s2 EtikettSet) (s3 EtikettSet) {
	s4 := MakeEtikettMutableSet()

	for _, e1 := range s2.Elements() {
		didAdd := false

		for _, e := range s1.Elements() {
			if strings.HasPrefix(e.String(), e1.String()) {
				didAdd = true
				s4.Add(e)
			}
		}

		if !didAdd {
			s4 = MakeEtikettMutableSet()
			return
		}
	}

	s3 = s4.Copy()

	return
}

func SubtractPrefix(s1 EtikettSet, e Etikett) (s2 EtikettSet) {
	s3 := MakeEtikettMutableSet()

	for _, e1 := range s1.Elements() {
		e2, _ := e1.LeftSubtract(e)

		if e2.String() == "" {
			continue
		}

		s3.Add(e2)
	}

	s2 = s3.Copy()

	return
}

func Description(s EtikettSet) string {
	sb := &strings.Builder{}
	first := true

	for _, e1 := range s.Sorted() {
		if !first {
			sb.WriteString(", ")
		}

		sb.WriteString(e1.String())

		first = false
	}

	return sb.String()
}
