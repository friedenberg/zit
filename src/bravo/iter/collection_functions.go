package iter

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

func AnyOrTrueEmpty[T any](c schnittstellen.Set[T], f func(T) bool) bool {
	if c.Len() == 0 {
		return true
	}

	return Any(c, f)
}

func AnyOrFalseEmpty[T any](c schnittstellen.Set[T], f func(T) bool) bool {
	if c.Len() == 0 {
		return false
	}

	return Any(c, f)
}

func Any[T any](c schnittstellen.Set[T], f func(T) bool) bool {
	err := c.Each(
		func(e T) (err error) {
			if f(e) {
				err = errTrue
			}

			return
		},
	)

	return IsErrTrue(err)
}

func All[T any](c schnittstellen.Set[T], f func(T) bool) bool {
	err := c.Each(
		func(e T) (err error) {
			if !f(e) {
				err = errFalse
			}

			return
		},
	)

	return !IsErrFalse(err)
}
