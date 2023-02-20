package values

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

func Equals[T schnittstellen.Equatable[T]](a T, b any) bool {
	b1, ok := b.(T)

	if !ok {
		return false
	}

	return a.Equals(b1)
}
