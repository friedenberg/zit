package iter

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

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
				err = errors.ErrTrue
			}

			return
		},
	)

	return errors.IsErrTrue(err)
}

func CheckAny[T any](c schnittstellen.Collection[T], f func(T) bool) bool {
	err := c.Each(
		func(e T) (err error) {
			if f(e) {
				err = errors.ErrTrue
			}

			return
		},
	)

	return errors.IsErrTrue(err)
}

func All[T any](c schnittstellen.Collection[T], f func(T) bool) bool {
	err := c.Each(
		func(e T) (err error) {
			if !f(e) {
				err = errors.ErrFalse
			}

			return
		},
	)

	return !errors.IsErrFalse(err)
}

func MakeFuncSetString[
	E any,
	EPtr schnittstellen.SetterPtr[E],
](
	c schnittstellen.Adder[E],
) schnittstellen.FuncSetString {
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

func Len(cs ...schnittstellen.Lenner) (n int) {
	for _, c := range cs {
		n += c.Len()
	}

	return
}

func Map[E schnittstellen.Value[E], F schnittstellen.Value[F]](
	in schnittstellen.SetLike[E],
	tr schnittstellen.FuncTransform[E, F],
	out schnittstellen.MutableSetLike[F],
) (err error) {
	if err = in.Each(
		func(e E) (err error) {
			var e1 F

			if e1, err = tr(e); err != nil {
				if IsStopIteration(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			return out.Add(e1)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func DerivedValues[E any, F any](
	c schnittstellen.SetLike[E],
	f schnittstellen.FuncTransform[E, F],
) (out []F, err error) {
	out = make([]F, 0, c.Len())

	if err = c.Each(
		func(e E) (err error) {
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

func DerivedValuesPtr[E any, EPtr schnittstellen.Ptr[E], F any](
	c schnittstellen.SetPtrLike[E, EPtr],
	f schnittstellen.FuncTransform[EPtr, F],
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
	wf schnittstellen.FuncIter[T],
) schnittstellen.FuncIter[T1] {
	return func(e T1) (err error) {
		if e1, ok := any(e).(T); ok {
			return wf(e1)
		}

		return
	}
}
