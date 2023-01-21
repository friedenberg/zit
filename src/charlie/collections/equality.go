package collections

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"golang.org/x/exp/constraints"
)

func EqualMapsEquatable[K constraints.Ordered, V schnittstellen.Equatable[V]](
	a, b map[K]V,
) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		v1, ok := b[k]

		if !ok {
			return false
		}

		if !v.Equals(v1) {
			return false
		}
	}

	return true
}

func EqualMapsOrdered[K constraints.Ordered, V constraints.Ordered](
	a, b map[K]V,
) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		v1, ok := b[k]

		if !ok {
			return false
		}

		if v != v1 {
			return false
		}
	}

	return true
}

func EqualSliceOrdered[V constraints.Ordered](
	a, b []V,
) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		v1 := b[i]

		if v != v1 {
			return false
		}
	}

	return true
}
