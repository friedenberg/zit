package kennung

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

func init() {
	collections.RegisterGob[Etikett]()
}

type EtikettSet = schnittstellen.Set[Etikett]

func MakeEtikettSet(es ...Etikett) (s EtikettSet) {
	return EtikettSet(collections.MakeSet((Etikett).String, es...))
}

func MakeSetStrings(vs ...string) (s EtikettSet, err error) {
	f := collections.MakeFlagCommasFromExisting(
		collections.SetterPolicyAppend,
		&s,
	)

	err = f.SetMany(vs...)

	return
}

func WithRemovedCommonPrefixes(s EtikettSet) (s2 EtikettSet) {
	es1 := collections.SortedValues(s)
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

	if len(exes) == 0 {
		exes = []Expander{ExpanderAll}
	}

	for _, e := range s.Elements() {
		e.Expanded(exes...).Each(s1.Add)
	}

	out = s1.ImmutableClone()

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

	s3 = s4.ImmutableClone()

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

	s2 = s3.ImmutableClone()

	return
}

func Description(s EtikettSet) string {
	return collections.StringCommaSeparated(s)
}
