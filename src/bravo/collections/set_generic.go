package collections

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Set[T any] struct {
	SetLike[T]
}

func MakeSetGeneric[T any](kf KeyFunc[T], es ...T) (out Set[T]) {
	out.SetLike = makeSetGeneric(kf, es...)

	return
}

func (s1 Set[T]) Copy() (out Set[T]) {
	s2 := makeSetGeneric[T](s1.Key)
	s2.open()
	defer s2.close()

	s1.Each(s2.add)

	out.SetLike = s2

	return
}

func (s1 Set[T]) MutableCopy() (s2 MutableSetLike[T]) {
	s2 = makeMutableSetGeneric[T](s1.Key)
	s1.Each(s2.Add)

	return
}

func (s Set[T]) WriterContainer() WriterFunc[T] {
	return func(e T) (err error) {
		k := s.Key(e)

		if k == "" {
			err = ErrEmptyKey[T]{Element: e}
			return
		}

		_, ok := s.Get(k)

		if !ok {
			err = io.EOF
		}

		return
	}
}

func WriterFuncNegate[T any](wf WriterFunc[T]) WriterFunc[T] {
	return func(e T) (err error) {
		err = wf(e)

		switch {
		case err == nil:
			err = io.EOF

		case errors.IsEOF(err):
			err = nil
		}

		return
	}
}

func (s1 Set[T]) Subtract(s2 Set[T]) (out Set[T]) {
	s3 := makeSetGeneric[T](s1.Key)
	s3.open()
	defer s3.close()

	s1.Chain(
		WriterFuncNegate(s2.WriterContainer()),
		s3.add,
	)

	out.SetLike = s3

	return
}

func (s1 Set[T]) Intersection(s2 SetLike[T]) (s3 MutableSetLike[T]) {
	s3 = MakeMutableSetGeneric[T](s1.Key)
	s22 := Set[T]{
		SetLike: s2,
	}

	s1.Chain(
		s22.WriterContainer(),
		s3.Add,
	)

	return
}

func (s1 Set[T]) Chain(fs ...WriterFunc[T]) error {
	return s1.Each(
		func(e T) (err error) {
			for _, f := range fs {
				if err = f(e); err != nil {
					if errors.IsEOF(err) {
						err = nil
					} else {
						err = errors.Wrap(err)
					}

					return
				}
			}

			return
		},
	)
}

func (s Set[T]) Elements() (out []T) {
	out = make([]T, s.Len())

	s.Each(
		func(e T) (err error) {
			out = append(out, e)
			return
		},
	)

	return
}

func (s Set[T]) Any() (e T) {
	s.Each(
		func(e1 T) (err error) {
			e = e1
			return io.EOF
		},
	)

	return
}

func (s Set[T]) All(f WriterFunc[T]) (ok bool) {
	err := s.Each(
		func(e T) (err error) {
			return f(e)
		},
	)

	return err == nil
}

func (a Set[T]) Equals(b Set[T]) (ok bool) {
	if a.Len() != b.Len() {
		return
	}

	ok = a.All(b.WriterContainer())

	return
}

func (outer Set[T]) ContainsSet(inner Set[T]) (ok bool) {
	if outer.Len() < inner.Len() {
		return
	}

	ok = inner.All(outer.WriterContainer())

	return
}
