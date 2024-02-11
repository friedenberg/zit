package values

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
)

func Equals[T schnittstellen.Equatable[T]](a T, b any) bool {
	{
		b1, ok := b.(T)

		if ok {
			return a.Equals(b1)
		}
	}

	return false
}

func EqualsPtr[T schnittstellen.Equatable[T], TPtr schnittstellen.Ptr[T]](
	a T,
	b any,
) bool {
	{
		b1, ok := b.(T)

		if ok {
			return a.Equals(b1)
		}
	}

	{
		b1, ok := b.(TPtr)

		if ok {
			return a.Equals(*b1)
		}
	}

	return false
}
