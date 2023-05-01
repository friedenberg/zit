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
	schnittstellen.MutableSet[*Transacted]
}

func MakeMutableSetUnique(c int) MutableSet {
	return MutableSet{
		MutableSet: collections.MakeMutableSet(
			func(sz *Transacted) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Sku.Kopf,
					sz.Sku.Schwanz,
					sz.Sku.TransactionIndex,
					sz.Sku.Kennung,
					sz.Sku.ObjekteSha,
				)
			},
		),
	}
}

func MakeMutableSetHinweis(c int) MutableSet {
	return MutableSet{
		MutableSet: collections.MakeMutableSet(
			func(sz *Transacted) string {
				if sz == nil {
					return ""
				}

				return collections.MakeKey(
					sz.Sku.Kennung,
				)
			},
		),
	}
}

func (s MutableSet) ToSliceHinweisen() (b []kennung.Hinweis) {
	b = make([]kennung.Hinweis, 0, s.Len())

	s.Each(
		func(z *Transacted) (err error) {
			b = append(b, z.Sku.Kennung)

			return
		},
	)

	return
}
