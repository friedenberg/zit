package collections

type setGenericAlias[T any] struct {
	Set[T]
}

type MutableSet[T any] struct {
	setGenericAlias[T]
	MutableSetLike[T]
}

func MakeMutableSetGeneric[T any](kf KeyFunc[T], es ...T) (out MutableSet[T]) {
	out.MutableSetLike = makeMutableSetGeneric(kf, es...)
	out.setGenericAlias = setGenericAlias[T]{Set: Set[T]{SetLike: out.MutableSetLike}}

	return
}
