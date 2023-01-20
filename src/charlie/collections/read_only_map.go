package collections

import "golang.org/x/exp/constraints"

type ReadOnlyMap[K constraints.Ordered, V any] map[K]V
