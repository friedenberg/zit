package iter

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
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
