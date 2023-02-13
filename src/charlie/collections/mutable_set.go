package collections

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type MutableSet[T any] struct {
	setAlias[T]
	schnittstellen.MutableSet[T]
}

func MakeMutableSet[T any](kf KeyFunc[T], es ...T) (out MutableSet[T]) {
	out.MutableSet = makeMutableSet(kf, es...)
	out.setAlias = setAlias[T]{Set: Set[T]{Set: out.MutableSet}}

	return
}

func (s MutableSet[T]) AddAndDoNotRepool(e T) (err error) {
	s.Add(e)
	err = ErrDoNotRepool
	return
}
