package iter

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func Elements[T any](s schnittstellen.Iterable[T]) (out []T) {
	out = make([]T, 0, s.Len())

	err := s.Each(
		func(v T) (err error) {
			out = append(out, v)
			return
		},
	)

	errors.PanicIfError(err)

	return
}
