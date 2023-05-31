package collections

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type delta[T schnittstellen.ValueLike] struct {
	Added, Removed schnittstellen.MutableSet[T]
}

func (d delta[T]) GetAdded() schnittstellen.Set[T] {
	return d.Added
}

func (d delta[T]) GetRemoved() schnittstellen.Set[T] {
	return d.Removed
}

func MakeSetDelta[T schnittstellen.ValueLike](
	from, to schnittstellen.Set[T],
) schnittstellen.Delta[T] {
	d := delta[T]{
		Added:   MakeMutableSetStringer[T](),
		Removed: from.MutableClone(),
	}

	to.Each(
		func(e T) (err error) {
			if from.Contains(e) {
				// had previously
			} else {
				// did not have previously
				d.Added.Add(e)
			}

			d.Removed.Del(e)

			return
		},
	)

	return d
}
