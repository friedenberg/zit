package etikett

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/proto_objekte"
)

type Set = proto_objekte.Set[Etikett, *Etikett]

func MakeSet(es ...Etikett) (s Set) {
	return Set(proto_objekte.MakeSet(es...))
}

func MakeSetStrings(vs ...string) (s Set, err error) {
	var s1 proto_objekte.Set[Etikett, *Etikett]

	if s1, err = proto_objekte.MakeSetStrings[Etikett, *Etikett](vs...); err != nil {
		err = errors.Wrap(err)
		return
	}

	s = Set(s1)

	return
}

// func (s *Set) Set(v string) (err error) {
// 	if s.closed {
// 		err = errors.Errorf("trying to mutate closed set")
// 		return
// 	}

// 	s.inner = make(map[string]Etikett, 1)

// 	es := strings.Split(v, ",")

// 	if len(es) == 0 {
// 		return
// 	}

// 	if es[0] == "" {
// 		return
// 	}

// 	for _, e := range es {
// 		var e1 Etikett

// 		if err = e1.Set(e); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}

// 		s.addOnlyExact(e1)
// 	}

// 	return
// }

func WithRemovedCommonPrefixes(s Set) (s2 Set) {
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

	s2 = MakeSet(es...)

	return
}

func Expanded(s Set, exes ...Expander) (out Set) {
	s1 := MakeMutableSet()

	for _, e := range s.Elements() {
		s1.Merge(e.Expanded(exes...))
	}

	out = s1.Copy()

	return
}

// func (s Set) String() string {
// 	return strings.Join(s.SortedString(), ", ")
// }

func IntersectPrefixes(s1 Set, s2 Set) (s3 Set) {
	s4 := MakeMutableSet()

	for _, e1 := range s2.Elements() {
		didAdd := false

		for _, e := range s1.Elements() {
			if strings.HasPrefix(e.String(), e1.String()) {
				didAdd = true
				s4.Add(e)
			}
		}

		if !didAdd {
			s4 = MakeMutableSet()
			return
		}
	}

	s3 = s4.Copy()

	return
}

func SubtractPrefix(s1 Set, e Etikett) (s2 Set) {
	s3 := MakeMutableSet()

	for _, e1 := range s1.Elements() {
		e2 := e1.LeftSubtract(e)

		if e2.String() == "" {
			continue
		}

		s3.Add(e2)
	}

	s2 = s3.Copy()

	return
}

func Description(s Set) string {
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
