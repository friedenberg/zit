package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/tridex"
)

type Indexed[T KennungLike[T], TPtr KennungLikePtr[T]] struct {
	Int            int
	Kennung        T
	SchwanzenCount int
	Count          int
	ExpandedAll    schnittstellen.Set[T]
	ExpandedRight  schnittstellen.Set[T]
	Tridex         schnittstellen.MutableTridex
}

func MakeIndexed[
	T KennungLike[T],
	TPtr KennungLikePtr[T],
](k T,
) (i *Indexed[T, TPtr]) {
	i = &Indexed[T, TPtr]{}
	i.ResetWithKennung(k)
	return
}

func (i *Indexed[T, TPtr]) ResetWithKennung(k T) {
	i.Kennung = k
	i.ExpandedAll = ExpandOne[T, TPtr](k, ExpanderAll)
	i.ExpandedRight = ExpandOne[T, TPtr](k, ExpanderRight)
	i.Tridex = tridex.Make(collections.SortedStrings[T](i.ExpandedRight)...)
}

func (z Indexed[T, TPtr]) GetInt() int {
	return 0
}

func (z Indexed[T, TPtr]) GetKennung() T {
	return z.Kennung
}

func (k Indexed[T, TPtr]) GetSchwanzenCount() int {
	return k.SchwanzenCount
}

func (k Indexed[T, TPtr]) GetCount() int {
	return k.Count
}

func (z Indexed[T, TPtr]) GetTridex() schnittstellen.Tridex {
	return z.Tridex
}

func (z Indexed[T, TPtr]) GetExpandedRight() schnittstellen.Set[T] {
	return z.ExpandedRight
}

func (z Indexed[T, TPtr]) GetExpandedAll() schnittstellen.Set[T] {
	return z.ExpandedAll
}

func (z *Indexed[T, TPtr]) Reset() {
	TPtr(&z.Kennung).Reset()
	z.SchwanzenCount = 0
	z.Count = 0
	z.ExpandedRight = collections.MakeSetStringer[T, TPtr]()
	z.ExpandedAll = collections.MakeSetStringer[T, TPtr]()
	z.Tridex = tridex.Make()
}
