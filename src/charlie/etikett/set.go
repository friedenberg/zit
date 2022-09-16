package etikett

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Set map[string]Etikett

func (s Set) Len() int {
	return len(s)
}

func MakeSet(es ...Etikett) (s Set) {
	s = make(Set)

	for _, e := range es {
		s.Add(e)
	}

	return
}

func NewSet(es ...Etikett) (s *Set) {
	s1 := MakeSet(es...)
	s = &s1

	return
}

func (es Set) AddString(v string) (err error) {
	var e Etikett

	if err = e.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	es.Add(e)

	return
}

func (es Set) Add(e Etikett) {
	expanded := e.Expanded(ExpanderRight{})
	intersection := es.Intersect(*expanded)
	es.Remove(intersection.Etiketten()...)

	es.addOnlyExact(e)
}

func (es Set) addOnlyExact(e Etikett) {
	es[e.String()] = e
}

func (s *Set) Set(v string) (err error) {
	es := strings.Split(v, ",")

	for _, e := range es {
		if err = s.AddString(e); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (es Set) Remove(es1 ...Etikett) {
	for _, e := range es1 {
		delete(es, e.String())
	}
}

func (es Set) RemovePrefixes(needle Etikett) {
	for haystack, _ := range es {
		if strings.HasPrefix(haystack, needle.String()) {
			delete(es, haystack)
		}
	}
}

func (a Set) Equals(b Set) bool {
	if len(a) != len(b) {
		return false
	}

	for ae, _ := range a {
		if _, ok := b[ae]; !ok {
			return false
		}
	}

	return true
}

func (s1 Set) Merge(s2 Set) {
	for _, e := range s2 {
		s1.Add(e)
	}
}

func (s1 Set) Copy() (s2 *Set) {
	s2 = NewSet()

	for _, e := range s1 {
		s2.Add(e)
	}

	return
}

func (s Set) Expanded(exes ...Expander) (s1 Set) {
	s1 = MakeSet()

	for _, e := range s {
		for _, e1 := range *e.Expanded(exes...) {
			s1.addOnlyExact(e1)
		}
	}

	return
}

func (s Set) String() string {
	return strings.Join(s.SortedString(), ", ")
}

func (s Set) Strings() (out []string) {
	out = make([]string, 0, len(s))

	for s, _ := range s {
		out = append(out, s)
	}

	return
}

func (es Set) Etiketten() (out []Etikett) {
	out = make([]Etikett, len(es))

	i := 0

	for _, e := range es {
		out[i] = e
		i++
	}

	return
}

func (es Set) Sorted() (out []Etikett) {
	out = es.Etiketten()

	sort.Slice(
		out,
		func(i, j int) bool {
			return out[i].String() < out[j].String()
		},
	)

	return
}

func (es Set) SortedString() (out []string) {
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

func (s Set) Contains(e Etikett) bool {
	return s.ContainsString(e.String())
}

func (s Set) ContainsString(es string) bool {
	_, ok := s[es]
	return ok
}

func (a Set) ContainsSet(b Set) bool {
	for _, e := range b {
		if !a.Contains(e) {
			return false
		}
	}

	return true
}

func (s1 Set) Subtract(s2 Set) (s3 Set) {
	s3 = MakeSet()

	for _, e1 := range s1 {
		if s2.Contains(e1) {
			continue
		}

		s3.Add(e1)
	}

	return
}

func (s1 Set) IntersectPrefixes(s2 Set) (s3 Set) {
	s3 = MakeSet()

	for _, e1 := range s2 {
		didAdd := false

		for _, e := range s1 {
			if strings.HasPrefix(e.String(), e1.String()) {
				didAdd = true
				s3.Add(e)
			}
		}

		if !didAdd {
			s3 = MakeSet()
			return
		}
	}

	return
}

func (s1 Set) Intersect(s2 Set) (s3 Set) {
	s3 = MakeSet()

	for _, e := range s1 {
		if s2.Contains(e) {
			s3.Add(e)
		}
	}

	return
}

func (s1 Set) SubtractPrefix(e Etikett) (s2 Set) {
	s2 = MakeSet()

	for _, e1 := range s1 {
		e2 := e1.LeftSubtract(e)

		if e2.String() == "" {
			continue
		}

		s2.Add(e2)
	}

	return
}

func (s Set) Any() (e Etikett) {
	for _, e1 := range s {
		e = e1
		break
	}

	return e
}

func (es Set) Description() string {
	sb := &strings.Builder{}
	first := true

	for _, e1 := range es.Sorted() {
		if !first {
			sb.WriteString(", ")
		}

		sb.WriteString(e1.String())

		first = false
	}

	return sb.String()
}

func (es Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(es.SortedString())
}

func (es *Set) UnmarshalJSON(b []byte) (err error) {
	es = NewSet()
	var vs []string

	if err = json.Unmarshal(b, &vs); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, v := range vs {
		if err = es.AddString(v); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
