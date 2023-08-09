package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/tridex"
)

type IndexedLike[
	T KennungLike[T],
	TPtr KennungLikePtr[T],
] struct {
	Int            int
	Kennung        T
	SchwanzenCount int
	Count          int
	ExpandedAll    schnittstellen.SetPtrLike[T, TPtr]
	ExpandedRight  schnittstellen.SetPtrLike[T, TPtr]
	Tridex         schnittstellen.MutableTridex
}

func MakeIndexed[
	T KennungLike[T],
	TPtr KennungLikePtr[T],
](k T,
) (i *IndexedLike[T, TPtr]) {
	i = &IndexedLike[T, TPtr]{}
	i.ResetWithKennung(k)
	return
}

func (i *IndexedLike[T, TPtr]) ResetWithKennung(k T) {
	i.Kennung = k
	i.ExpandedAll = ExpandOne[T, TPtr](&k, ExpanderAll)
	i.ExpandedRight = ExpandOne[T, TPtr](&k, ExpanderRight)
	i.Tridex = tridex.Make(iter.SortedStrings[T](i.ExpandedRight)...)
}

func (z IndexedLike[T, TPtr]) GetInt() int {
	return 0
}

func (z IndexedLike[T, TPtr]) GetKennung() T {
	return z.Kennung
}

func (k IndexedLike[T, TPtr]) GetSchwanzenCount() int {
	return k.SchwanzenCount
}

func (k IndexedLike[T, TPtr]) GetCount() int {
	return k.Count
}

func (z IndexedLike[T, TPtr]) GetTridex() schnittstellen.Tridex {
	return z.Tridex
}

func (z IndexedLike[T, TPtr]) GetExpandedRight() schnittstellen.SetPtrLike[T, TPtr] {
	return z.ExpandedRight
}

func (z IndexedLike[T, TPtr]) GetExpandedAll() schnittstellen.SetPtrLike[T, TPtr] {
	return z.ExpandedAll
}

func (z *IndexedLike[T, TPtr]) Reset() {
	TPtr(&z.Kennung).Reset()
	z.SchwanzenCount = 0
	z.Count = 0
	z.ExpandedRight = collections_ptr.MakeMutableValueSetValue[T, TPtr](nil)
	z.ExpandedAll = collections_ptr.MakeMutableValueSetValue[T, TPtr](nil)
	z.Tridex = tridex.Make()
}
