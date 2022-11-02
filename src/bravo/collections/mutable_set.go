package collections

type MutableSet[T any] struct {
	setAlias[T]
	MutableSetLike[T]
}

func MakeMutableSetGeneric[T any](kf KeyFunc[T], es ...T) (out MutableSet[T]) {
	out.MutableSetLike = makeMutableSet(kf, es...)
	out.setAlias = setAlias[T]{Set: Set[T]{SetLike: out.MutableSetLike}}

	return
}
