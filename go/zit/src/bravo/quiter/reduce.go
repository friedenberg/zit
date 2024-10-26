package quiter

import (
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func Elements[T any](s interfaces.Iterable[T]) (out []T) {
	if s == nil {
		return
	}

	out = make([]T, 0, s.Len())

	for v := range s.All() {
		out = append(out, v)
	}

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

	for v := range s.All() {
		out = append(out, v)
	}

	sort.Slice(out, func(i, j int) bool {
		return sf(out[i], out[j])
	})

	return
}
