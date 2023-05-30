package kennung_index

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/tridex"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type indexed[T kennung.KennungLike[T], TPtr kennung.KennungLikePtr[T]] struct {
	Kennung        T
	SchwanzenCount int
	Count          int
	ExpandedAll    schnittstellen.Set[T]
	ExpandedRight  schnittstellen.Set[T]
	Tridex         schnittstellen.MutableTridex
}

func makeIndexed[
	T kennung.KennungLike[T],
	TPtr kennung.KennungLikePtr[T],
](k T,
) (i *indexed[T, TPtr]) {
	i = &indexed[T, TPtr]{}
	i.ResetWithKennung(k)
	return
}

func (i *indexed[T, TPtr]) ResetWithKennung(k T) {
	i.Kennung = k
	i.ExpandedAll = kennung.ExpandOne[T, TPtr](k, kennung.ExpanderAll)
	i.ExpandedRight = kennung.ExpandOne[T, TPtr](k, kennung.ExpanderRight)
	i.Tridex = tridex.Make(collections.SortedStrings[T](i.ExpandedRight)...)
}

func (z indexed[T, TPtr]) GetKennung() T {
	return z.Kennung
}

func (k indexed[T, TPtr]) GetSchwanzenCount() int {
	return k.SchwanzenCount
}

func (k indexed[T, TPtr]) GetCount() int {
	return k.Count
}

func (z indexed[T, TPtr]) GetTridex() schnittstellen.Tridex {
	return z.Tridex
}

func (z indexed[T, TPtr]) GetExpandedRight() schnittstellen.Set[T] {
	return z.ExpandedRight
}

func (z indexed[T, TPtr]) GetExpandedAll() schnittstellen.Set[T] {
	return z.ExpandedAll
}

func (z *indexed[T, TPtr]) Reset() {
	TPtr(&z.Kennung).Reset()
	z.SchwanzenCount = 0
	z.Count = 0
	z.ExpandedRight = collections.MakeSetStringer[T]()
	z.ExpandedAll = collections.MakeSetStringer[T]()
	z.Tridex = tridex.Make()
}
