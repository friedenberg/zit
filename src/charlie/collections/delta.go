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
	s1, s2 schnittstellen.Set[T],
) schnittstellen.Delta[T] {
	d := delta[T]{
		Added:   MakeMutableSetStringer[T](),
		Removed: s1.MutableClone(),
	}

	s2.Each(
		func(e T) (err error) {
			if s1.Contains(e) {
				// zettel had etikett previously
			} else {
				// zettel did not have etikett previously
				d.Added.Add(e)
			}

			d.Removed.Del(e)

			return
		},
	)

	return d
}
