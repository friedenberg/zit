package iter

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

func AddOrReplaceIfGreaterCustom[T interface {
	schnittstellen.Stringer
	schnittstellen.ValueLike
	schnittstellen.Lessor[T]
}](c schnittstellen.MutableSetLike[T], b T, kf func(T) string) (err error) {
	a, ok := c.Get(kf(b))

	if !ok || a.Less(b) {
		return c.AddCustomKey(b, kf)
	}

	return
}

func AddOrReplaceIfGreater[T interface {
	schnittstellen.Stringer
	schnittstellen.ValueLike
	schnittstellen.Lessor[T]
}](c schnittstellen.MutableSetLike[T], b T) (err error) {
	a, ok := c.Get(b.String())

	if !ok || a.Less(b) {
		return c.Add(b)
	}

	return
}

func AnyOrTrueEmpty[T any](
	c schnittstellen.Collection[T],
	f func(T) bool,
) bool {
	if c.Len() == 0 {
		return true
	}

	return Any(c, f)
}

func AnyOrFalseEmpty[T any](
	c schnittstellen.Collection[T],
	f func(T) bool,
) bool {
	if c.Len() == 0 {
		return false
	}

	return Any(c, f)
}

func Any[T any](c schnittstellen.Collection[T], f func(T) bool) bool {
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

func All[T any](c schnittstellen.Collection[T], f func(T) bool) bool {
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
