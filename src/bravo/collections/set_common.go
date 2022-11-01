package collections

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func Any[T any](s SetLike[T]) (e T) {
	s.Each(
		func(e1 T) (err error) {
			e = e1
			return io.EOF
		},
	)

	return
}

func Chain[T any](s SetLike[T], fs ...WriterFunc[T]) error {
	return s.Each(
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

func Elements[T any](s SetLike[T]) (out []T) {
	out = make([]T, s.Len())

	s.Each(
		func(e T) (err error) {
			out = append(out, e)
			return
		},
	)

	return
}

func All[T any](s SetLike[T], f WriterFunc[T]) (ok bool) {
	err := s.Each(
		func(e T) (err error) {
			return f(e)
		},
	)

	return err == nil
}

func Equals[T any](a, b SetLike[T]) (ok bool) {
	if a.Len() != b.Len() {
		return
	}

	ok = All(
		a,
		func(e T) (err error) {
			if !b.Contains(e) {
				err = io.EOF
				return
			}

			return
		},
	)

	return
}

func ContainsSet[T any](outer, inner SetLike[T]) (ok bool) {
	if outer.Len() < inner.Len() {
		return
	}

	ok = All(
		inner,
		func(e T) (err error) {
			if !outer.Contains(e) {
				err = io.EOF
				return
			}

			return
		},
	)

	return
}

func Intersection[T any](s1, s2 SetLike[T]) (s3 MutableSetLike[T]) {
	s3 = MakeMutableSetGeneric[T](s1.Key)

	Chain(
		s2,
		s1.WriterContainer(),
		s3.Add,
	)

	return
}
