package etikett

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Set struct {
  closed bool
	inner map[string]Etikett
}

func (s Set) Len() int {
	return len(s.inner)
}

func MakeSet(es ...Etikett) (s Set) {
	s.inner = make(map[string]Etikett, len(es))

	for _, e := range es {
		s.addOnlyExact(e)
	}

  s.closed = true

	return
}

func MakeSetFromStrings(es ...string) (s Set, err error) {
	s.inner = make(map[string]Etikett, len(es))

	for _, v := range es {
    var e Etikett

    if err = e.Set(v); err != nil {
      err = errors.Wrap(err)
      return
    }

		s.addOnlyExact(e)
	}

  s.closed = true

	return
}

func (es *Set) addOnlyExact(e Etikett) {
	es.inner[e.String()] = e
}

func (s *Set) Set(v string) (err error) {
  if s.closed {
    err = errors.Errorf("trying to mutate closed set")
    return
  }

	s.inner = make(map[string]Etikett, 1)

	es := strings.Split(v, ",")

  if len(es) == 0 {
    return
  }

  if es[0] == "" {
    return
  }

	for _, e := range es {
    var e1 Etikett

		if err = e1.Set(e); err != nil {
			err = errors.Wrap(err)
			return
		}

    s.addOnlyExact(e1)
	}

	return
}

func (s Set) WithRemovedCommonPrefixes() (s2 Set) {
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

func (es *Set) Remove(es1 ...Etikett) {
	for _, e := range es1 {
		delete(es.inner, e.String())
	}
}

func (es *Set) RemovePrefixes(needle Etikett) {
	for haystack, _ := range es.inner {
		if strings.HasPrefix(haystack, needle.String()) {
			delete(es.inner, haystack)
		}
	}
}

func (a Set) Equals(b Set) bool {
	if len(a.inner) != len(b.inner) {
		return false
	}

	for ae, _ := range a.inner {
		if _, ok := b.inner[ae]; !ok {
			return false
		}
	}

	return true
}

func (s1 *Set) Merge(s2 Set) {
	for _, e := range s2.inner {
		s1.addOnlyExact(e)
	}
}

func (s1 Set) Copy() (s2 Set) {
	s2 = MakeSet()

	for _, e := range s1.inner {
		s2.addOnlyExact(e)
	}

	return
}

func (s1 Set) MutableCopy() (s2 MutableSet) {
	s2 = MakeMutableSet()

	for _, e := range s1.inner {
		s2.addOnlyExact(e)
	}

	return
}

func (s Set) Expanded(exes ...Expander) (s1 Set) {
	s1 = MakeSet()

	for _, e := range s.inner {
		for _, e1 := range e.Expanded(exes...).inner {
			s1.addOnlyExact(e1)
		}
	}

	return
}

func (s Set) String() string {
	return strings.Join(s.SortedString(), ", ")
}

func (s Set) Strings() (out []string) {
	out = make([]string, 0, len(s.inner))

	for s, _ := range s.inner {
		out = append(out, s)
	}

	return
}

func (es Set) Etiketten() (out []Etikett) {
	out = make([]Etikett, len(es.inner))

	i := 0

	for _, e := range es.inner {
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
	out = make([]string, len(es.inner))

	i := 0

	for _, e := range es.inner {
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
	_, ok := s.inner[es]
	return ok
}

func (a Set) ContainsSet(b Set) bool {
	for _, e := range b.inner {
		if !a.Contains(e) {
			return false
		}
	}

	return true
}

func (s1 Set) Subtract(s2 Set) (s3 Set) {
	s3 = MakeSet()

	for _, e1 := range s1.inner {
		if s2.Contains(e1) {
			continue
		}

		s3.addOnlyExact(e1)
	}

	return
}

func (s1 Set) IntersectPrefixes(s2 Set) (s3 Set) {
	s3 = MakeSet()

	for _, e1 := range s2.inner {
		didAdd := false

		for _, e := range s1.inner {
			if strings.HasPrefix(e.String(), e1.String()) {
				didAdd = true
				s3.addOnlyExact(e)
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

	for _, e := range s1.inner {
		if s2.Contains(e) {
			s3.addOnlyExact(e)
		}
	}

	return
}

func (s1 *Set) Withdraw(e Etikett) (s2 Set) {
	s2 = MakeSet()

	for _, e1 := range s1.inner {
		if e1.Contains(e) {
			s2.addOnlyExact(e1)
		}
	}

	s1.Remove(s2.Etiketten()...)

	return
}

func (s1 Set) SubtractPrefix(e Etikett) (s2 Set) {
	s2 = MakeSet()

	for _, e1 := range s1.inner {
		e2 := e1.LeftSubtract(e)

		if e2.String() == "" {
			continue
		}

		s2.addOnlyExact(e2)
	}

	return
}

func (s Set) Any() (e Etikett) {
	for _, e1 := range s.inner {
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
	*es = MakeSet()
	var vs []string

	if err = json.Unmarshal(b, &vs); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, v := range vs {
    var e Etikett

		if err = e.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

    es.addOnlyExact(e)
	}

	return
}

func (s Set) MarshalBinary() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *Set) UnmarshalBinary(text []byte) (err error) {
	s.inner = make(map[string]Etikett)

	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
