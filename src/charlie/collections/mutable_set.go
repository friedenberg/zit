package collections

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type MutableSet[T any] struct {
	setAlias[T]
	schnittstellen.MutableSetLike[T]
}

func MakeMutableSet[T any](kf KeyFunc[T], es ...T) (out MutableSet[T]) {
	out.MutableSetLike = makeMutableSet(kf, es...)
	out.setAlias = setAlias[T]{Set: Set[T]{SetLike: out.MutableSetLike}}

	return
}

func (s MutableSet[T]) AddAndDoNotRepool(e T) (err error) {
	s.Add(e)
	err = ErrDoNotRepool
	return
}
