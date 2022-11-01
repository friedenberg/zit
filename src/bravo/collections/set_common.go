package collections

import (
	"io"
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
