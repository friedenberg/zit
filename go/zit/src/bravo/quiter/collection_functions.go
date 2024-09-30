package quiter

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func CheckAnyOrTrueEmpty[T any](
	c interfaces.Collection[T],
	f func(T) bool,
) bool {
	if c.Len() == 0 {
		return true
	}

	return CheckAny(c, f)
}

func CheckAnyOrFalseEmpty[T any](
	c interfaces.Collection[T],
	f func(T) bool,
) bool {
	if c.Len() == 0 {
		return false
	}

	return CheckAny(c, f)
}

func CheckAnyPtr[
	T any,
	TPtr interfaces.Ptr[T],
](
	c interfaces.CollectionPtr[T, TPtr],
	f func(TPtr) bool,
) bool {
	err := c.EachPtr(
		func(e TPtr) (err error) {
			if f(e) {
				err = errors.ErrTrue
			}

			return
		},
	)

	return errors.IsErrTrue(err)
}

func CheckAny[T any](c interfaces.Collection[T], f func(T) bool) bool {
	for e := range c.All() {
		if f(e) {
			return true
		}
	}

	return false
}

func All[T any](c interfaces.Collection[T], f func(T) bool) bool {
	for e := range c.All() {
		if !f(e) {
			return false
		}
	}

	return true
}

func MakeFuncSetString[
	E any,
	EPtr interfaces.SetterPtr[E],
](
	c interfaces.Adder[E],
) interfaces.FuncSetString {
	return func(v string) (err error) {
		return AddString[E, EPtr](c, v)
	}
}

// func ContainsKey(
// 	id schnittstellen.Stringer,
// 	cs ...schnittstellen.ContainsKeyer,
// ) (ok bool) {
// 	for _, c := range cs {
// 		if c.ContainsKey(id.String()) {
// 			return true
// 		}
// 	}

// 	return false
// }

func Len(cs ...interfaces.Lenner) (n int) {
	for _, c := range cs {
		n += c.Len()
	}

	return
}

func DerivedValues[E any, F any](
	c interfaces.SetLike[E],
	f interfaces.FuncTransform[E, F],
) (out []F, err error) {
	out = make([]F, 0, c.Len())

	for e := range c.All() {
		var e1 F

		if e1, err = f(e); err != nil {
			if IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		out = append(out, e1)
	}

	return
}

func DerivedValuesPtr[E any, EPtr interfaces.Ptr[E], F any](
	c interfaces.SetPtrLike[E, EPtr],
	f interfaces.FuncTransform[EPtr, F],
) (out []F, err error) {
	out = make([]F, 0, c.Len())

	if err = c.EachPtr(
		func(e EPtr) (err error) {
			var e1 F

			if e1, err = f(e); err != nil {
				if IsStopIteration(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			out = append(out, e1)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeFuncTransformer[T any, T1 any](
	wf interfaces.FuncIter[T],
) interfaces.FuncIter[T1] {
	return func(e T1) (err error) {
		if e1, ok := any(e).(T); ok {
			return wf(e1)
		}

		return
	}
}
