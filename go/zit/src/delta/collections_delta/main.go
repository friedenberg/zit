package collections_delta

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
)

type delta[T schnittstellen.ValueLike] struct {
	Added, Removed schnittstellen.MutableSetLike[T]
}

func (d delta[T]) GetAdded() schnittstellen.SetLike[T] {
	return d.Added
}

func (d delta[T]) GetRemoved() schnittstellen.SetLike[T] {
	return d.Removed
}

func MakeSetDelta[T schnittstellen.ValueLike](
	from, to schnittstellen.SetLike[T],
) schnittstellen.Delta[T] {
	d := delta[T]{
		Added:   collections_value.MakeMutableValueSet[T](nil),
		Removed: from.CloneMutableSetLike(),
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
