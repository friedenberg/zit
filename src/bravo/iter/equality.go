package iter

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

func SetEquals[T any](
	a, b schnittstellen.SetLike[T],
) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	if a.Len() != b.Len() {
		return false
	}

	err := a.Each(
		func(e T) (err error) {
			if ok := b.Contains(e); !ok {
				err = errFalse
				return
			}

			return
		},
	)

	if err != nil {
		return false
	}

	return true
}

func SetEqualsPtr[T any, TPtr schnittstellen.Ptr[T]](
	a, b schnittstellen.SetPtrLike[T, TPtr],
) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	if a.Len() != b.Len() {
		return false
	}

	err := a.EachPtr(
		func(e TPtr) (err error) {
			k := b.KeyPtr(e)

			if ok := b.ContainsKey(k); !ok {
				err = errFalse
				return
			}

			return
		},
	)
	if err != nil {
		return false
	}

	return true
}
