package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type HeapTransacted = collections.Heap[Transacted, *Transacted]

func MakeHeapTransacted() HeapTransacted {
	return collections.MakeHeap[Transacted, *Transacted]()
}

type MutableSet struct {
	schnittstellen.MutableSetLike[*Transacted]
}

func MakeMutableSetUnique(c int) MutableSet {
	return MutableSet{
		MutableSetLike: collections.MakeMutableSet(
			func(sz *Transacted) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Sku.Kopf,
					sz.Sku.GetTai(),
					sz.Sku.TransactionIndex,
					sz.Sku.GetKennung(),
					sz.Sku.ObjekteSha,
				)
			},
		),
	}
}

func MakeMutableSetHinweis(c int) MutableSet {
	return MutableSet{
		MutableSetLike: collections.MakeMutableSet(
			func(sz *Transacted) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Sku.GetKennung(),
				)
			},
		),
	}
}

func (s MutableSet) ToSliceHinweisen() (b []kennung.Hinweis) {
	b = make([]kennung.Hinweis, 0, s.Len())

	s.Each(
		func(z *Transacted) (err error) {
			b = append(b, z.Sku.GetKennung())

			return
		},
	)

	return
}
