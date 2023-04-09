package schnittstellen

type ImmutableCloner[T any] interface {
	ImmutableClone() T
}

type MutableCloner[T any] interface {
	MutableClone() T
}

// type UnionImmutableCloner[A interface{}, B any] interface {
// 	A
// 	ImmutableClone[B]
// }

func ImmutableClone[T ImmutableCloner[T]](o T) T {
	return o.ImmutableClone()
}
