package iter

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

func AddOrReplaceIfGreater[T interface {
	schnittstellen.Stringer
	schnittstellen.ValueLike
	schnittstellen.Lessor[T]
}](c schnittstellen.MutableSetLike[T], b T) (err error) {
	a, ok := c.Get(c.Key(b))

	if !ok || a.Less(b) {
		return c.Add(b)
	}

	return
}

func CheckAnyOrTrueEmpty[T any](
	c schnittstellen.Collection[T],
	f func(T) bool,
) bool {
	if c.Len() == 0 {
		return true
	}

	return CheckAny(c, f)
}

func CheckAnyOrFalseEmpty[T any](
	c schnittstellen.Collection[T],
	f func(T) bool,
) bool {
	if c.Len() == 0 {
		return false
	}

	return CheckAny(c, f)
}

func CheckAnyPtr[
	T any,
	TPtr schnittstellen.Ptr[T],
](
	c schnittstellen.CollectionPtr[T, TPtr],
	f func(TPtr) bool,
) bool {
	err := c.EachPtr(
		func(e TPtr) (err error) {
			if f(e) {
				err = errTrue
			}

			return
		},
	)

	return IsErrTrue(err)
}

func CheckAny[T any](c schnittstellen.Collection[T], f func(T) bool) bool {
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
