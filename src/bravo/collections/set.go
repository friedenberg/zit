package collections

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Set[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}] struct {
	keyFunc func(T) string
	closed  bool
	inner   map[string]T
}

func MakeSet[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](es ...T) (s Set[T, T1]) {
	s.inner = make(map[string]T, len(es))
	s.open()
	defer s.close()

	for _, e := range es {
		s.add(e)
	}

	return
}

func MakeSetStrings[T ProtoObjekte, T1 interface {
	*T
	ProtoObjektePointer
}](es ...string) (s Set[T, T1], err error) {
	s.inner = make(map[string]T, len(es))
	s.open()
	defer s.close()

	for _, e := range es {
		e1 := T1(new(T))

		if err = e1.Set(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		s.add(*e1)
	}

	return
}

func (s *Set[T, T1]) open() {
	s.closed = false
}

func (s *Set[T, T1]) close() {
	s.closed = true
}

func (s Set[T, T1]) Len() int {
	return len(s.inner)
}

func (s Set[T, T1]) GetKeyFunc() func(T) string {
	if s.keyFunc == nil {
		return func(e T) string {
			return e.String()
		}
	}

	return s.keyFunc
}

func (es Set[T, T1]) add(e T) {
	if es.closed {
		panic("trying to add etikett to closed set")
	}

	es.inner[es.GetKeyFunc()(e)] = e
}

func (s Set[T, T1]) Each(f func(T)) {
	for _, v := range s.inner {
		f(v)
	}
}

func (s *Set[T, T1]) Set(v string) (err error) {
	if s == nil {
		s1 := MakeSet[T, T1]()
		s = &s1
	}

	if s.closed {
		err = errors.Errorf("trying to mutate closed set")
		return
	}

	s.inner = make(map[string]T, 1)

	es := strings.Split(v, ",")

	if len(es) == 0 {
		return
	}

	if es[0] == "" {
		return
	}

	for _, e := range es {
		e1 := T1(new(T))

		if err = e1.Set(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		s.add(*e1)
	}

	return
}

func (es *Set[T, T1]) Remove(es1 ...T) {
	for _, e := range es1 {
		delete(es.inner, e.String())
	}
}

func (es *Set[T, T1]) RemovePrefixes(needle T) {
	for haystack, _ := range es.inner {
		if strings.HasPrefix(haystack, needle.String()) {
			delete(es.inner, haystack)
		}
	}
}

func (a Set[T, T1]) Equals(b Set[T, T1]) bool {
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

func (s1 Set[T, T1]) Copy() (s2 Set[T, T1]) {
	s2 = MakeSet[T, T1]()
	s2.open()
	defer s2.close()

	for _, e := range s1.inner {
		s2.add(e)
	}

	return
}

func (s1 Set[T, T1]) MutableCopy() (s2 MutableSet[T, T1]) {
	s2 = MakeMutableSet[T, T1]()

	for _, e := range s1.inner {
		s2.Add(e)
	}

	return
}

func (s Set[T, T1]) Strings() (out []string) {
	out = make([]string, 0, len(s.inner))

	for s, _ := range s.inner {
		out = append(out, s)
	}

	return
}

func (es Set[T, T1]) Elements() (out []T) {
	out = make([]T, len(es.inner))

	i := 0

	for _, e := range es.inner {
		out[i] = e
		i++
	}

	return
}

func (es Set[T, T1]) Sorted() (out []T) {
	out = es.Elements()

	sort.Slice(
		out,
		func(i, j int) bool {
			return out[i].String() < out[j].String()
		},
	)

	return
}

func (es Set[T, T1]) SortedString() (out []string) {
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

func (s Set[T, T1]) Contains(e T) bool {
	return s.ContainsString(s.GetKeyFunc()(e))
}

func (s Set[T, T1]) ContainsString(es string) bool {
	_, ok := s.inner[es]
	return ok
}

func (a Set[T, T1]) ContainsSet(b Set[T, T1]) bool {
	for _, e := range b.inner {
		if !a.Contains(e) {
			return false
		}
	}

	return true
}

func (s1 Set[T, T1]) Subtract(s2 Set[T, T1]) (s3 Set[T, T1]) {
	s3 = MakeSet[T, T1]()

	for _, e1 := range s1.inner {
		if s2.Contains(e1) {
			continue
		}

		s3.add(e1)
	}

	return
}

func (s1 Set[T, T1]) IntersectPrefixes(s2 Set[T, T1]) (s3 Set[T, T1]) {
	s3 = MakeSet[T, T1]()
	s3.open()
	defer s3.close()

	for _, e1 := range s2.inner {
		didAdd := false

		for _, e := range s1.inner {
			if strings.HasPrefix(e.String(), e1.String()) {
				didAdd = true
				s3.add(e)
			}
		}

		if !didAdd {
			s3 = MakeSet[T, T1]()
			return
		}
	}

	return
}

func (s1 Set[T, T1]) Intersect(s2 Set[T, T1]) (s3 Set[T, T1]) {
	s3 = MakeSet[T, T1]()

	for _, e := range s1.inner {
		if s2.Contains(e) {
			s3.add(e)
		}
	}

	return
}

func (s Set[T, T1]) Any() (e T) {
	for _, e1 := range s.inner {
		e = e1
		break
	}

	return e
}

func (es Set[T, T1]) Description() string {
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

func (s Set[T, T1]) String() string {
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

func (es Set[T, T1]) MarshalJSON() ([]byte, error) {
	return json.Marshal(es.SortedString())
}

func (es *Set[T, T1]) UnmarshalJSON(b []byte) (err error) {
	*es = MakeSet[T, T1]()

	es.open()
	defer es.close()

	var vs []string

	if err = json.Unmarshal(b, &vs); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, v := range vs {
		var e T1

		if err = e.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		es.add(*e)
	}

	return
}

func (s Set[T, T1]) MarshalBinary() (text []byte, err error) {
	text = []byte(s.String())

	return
}

func (s *Set[T, T1]) UnmarshalBinary(text []byte) (err error) {
	s.inner = make(map[string]T)

	s.open()
	defer s.close()

	if err = s.Set(string(text)); err != nil {
		return
	}

	return
}
