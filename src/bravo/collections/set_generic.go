package collections

import (
	"io"
	"reflect"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type SetGeneric[T any] struct {
	keyFunc func(T) string
	closed  bool
	inner   map[string]T
}

func MakeSetGeneric[T any](kf KeyFunc[T], es ...T) (s SetGeneric[T]) {
	t := *new(T)

	switch reflect.TypeOf(t).Kind() {
	// case reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
	case reflect.Ptr:
		kf(t)
	}

	s.keyFunc = kf
	s.inner = make(map[string]T, len(es))
	s.open()
	defer s.close()

	for _, e := range es {
		s.add(e)
	}

	return
}

func (s *SetGeneric[T]) open() {
	s.closed = false
}

func (s *SetGeneric[T]) close() {
	s.closed = true
}

func (s SetGeneric[T]) Len() int {
	return len(s.inner)
}

func (s SetGeneric[T]) Key(e T) string {
	return s.keyFunc(e)
}

func (s SetGeneric[T]) Get(k string) (e T, ok bool) {
	e, ok = s.inner[k]
	return
}

func (s SetGeneric[T]) ContainsKey(k string) (ok bool) {
	if k == "" {
		return
	}

	_, ok = s.inner[k]

	return
}

func (s SetGeneric[T]) Contains(e T) (ok bool) {
	return s.ContainsKey(s.Key(e))
}

func (es SetGeneric[T]) add(e T) {
	if es.closed {
		panic("trying to add etikett to closed set")
	}

	es.inner[es.Key(e)] = e
}

// TODO should this be locked mutable writes?
func (s SetGeneric[T]) Each(wf WriterFunc[T]) (err error) {
	for _, v := range s.inner {
		if err = wf(v); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (s1 SetGeneric[T]) Copy() (s2 SetGeneric[T]) {
	s2 = MakeSetGeneric[T](s1.Key)
	s2.open()
	defer s2.close()

	for _, e := range s1.inner {
		s2.add(e)
	}

	return
}

// func (s1 SetGeneric[T]) MutableCopy() (s2 MutableSet[T]) {
// 	s2 = MakeMutableSet[T]()

// 	for _, e := range s1.inner {
// 		s2.Add(e)
// 	}

// 	return
// }

// func (es SetGeneric[T]) Sorted() (out []T) {
// 	out = es.Elements()

// 	sort.Slice(
// 		out,
// 		func(i, j int) bool {
// 			return out[i].String() < out[j].String()
// 		},
// 	)

// 	return
// }

// func (es SetGeneric[T]) SortedString() (out []string) {
// 	out = make([]string, len(es.inner))

// 	i := 0

// 	for _, e := range es.inner {
// 		out[i] = e.String()
// 		i++
// 	}

// 	sort.Slice(
// 		out,
// 		func(i, j int) bool {
// 			return out[i] < out[j]
// 		},
// 	)

// 	return
// }

func (s MutableSetGeneric[T]) WriterContainer() WriterFunc[T] {
	return func(e T) (err error) {
		k := s.Key(e)

		if k == "" {
			err = ErrEmptyKey[T]{Element: e}
			return
		}

		_, ok := s.inner[k]

		if !ok {
			err = io.EOF
		}

		return
	}
}

// func (s SetGeneric[T]) ContainsString(es string) bool {
// 	_, ok := s.inner[es]
// 	return ok
// }

func (s1 SetGeneric[T]) Subtract(s2 SetGeneric[T]) (s3 SetGeneric[T]) {
	s3 = MakeSetGeneric[T](s1.Key)

	for _, e1 := range s1.inner {
		if s2.Contains(e1) {
			continue
		}

		s3.add(e1)
	}

	return
}

// func (s1 SetGeneric[T]) IntersectPrefixes(s2 SetGeneric[T]) (s3 SetGeneric[T]) {
// 	s3 = MakeSetGeneric[T](s1.KeyFunc())
// 	s3.open()
// 	defer s3.close()

// 	for _, e1 := range s2.inner {
// 		didAdd := false

// 		for _, e := range s1.inner {
// 			if strings.HasPrefix(e.String(), e1.String()) {
// 				didAdd = true
// 				s3.add(e)
// 			}
// 		}

// 		if !didAdd {
// 			s3 = MakeSetGeneric[T](s1.KeyFunc())
// 			return
// 		}
// 	}

// 	return
// }

func (s1 SetGeneric[T]) Intersect(s2 SetGeneric[T]) (s3 SetGeneric[T]) {
	s3 = MakeSetGeneric[T](s1.Key)

	for _, e := range s1.inner {
		if s2.Contains(e) {
			s3.add(e)
		}
	}

	return
}

func (s SetGeneric[T]) Any() (e T) {
	for _, e1 := range s.inner {
		e = e1
		break
	}

	return e
}
