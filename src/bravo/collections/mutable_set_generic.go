package collections

type setGenericAlias[T any] struct {
	SetGeneric[T]
}

type MutableSetGeneric[T any] struct {
	setGenericAlias[T]
	MutableSetLike[T]
}

func MakeMutableSetGeneric[T any](kf KeyFunc[T], es ...T) (out MutableSetGeneric[T]) {
	out.MutableSetLike = makeMutableSetGeneric(kf, es...)
	out.setGenericAlias = setGenericAlias[T]{SetGeneric: SetGeneric[T]{SetLike: out.MutableSetLike}}

	return
}
