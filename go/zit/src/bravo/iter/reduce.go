package iter

import (
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func Elements[T any](s interfaces.Iterable[T]) (out []T) {
	if s == nil {
		return
	}

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

func ElementsSorted[T any](
	s interfaces.Iterable[T],
	sf func(T, T) bool,
) (out []T) {
	if s == nil {
		return
	}

	out = make([]T, 0, s.Len())

	err := s.Each(
		func(v T) (err error) {
			out = append(out, v)
			return
		},
	)

	errors.PanicIfError(err)

	sort.Slice(out, func(i, j int) bool {
		return sf(out[i], out[j])
	})

	return
}
