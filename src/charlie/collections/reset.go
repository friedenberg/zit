package collections

import "golang.org/x/exp/constraints"

func ResetMap[K constraints.Ordered, V any](in map[K]V) (out map[K]V) {
	if in == nil {
		out = make(map[K]V)
	} else {
		for k := range in {
			delete(in, k)
		}

		out = in
	}

	return
}

func ResetSlice[V any](in []V) (out []V) {
	if in == nil {
		out = make([]V, 0)
	} else {
		out = in[:0]
	}

	return
}
