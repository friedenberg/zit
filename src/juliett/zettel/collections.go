package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
)

type HeapTransacted = collections.Heap[sku.TransactedZettel, *sku.TransactedZettel]

func MakeHeapTransacted() HeapTransacted {
	return collections.MakeHeap[sku.TransactedZettel, *sku.TransactedZettel]()
}

type (
	MutableSet = schnittstellen.MutableSetLike[*sku.TransactedZettel]
)

type TransactedUniqueKeyer struct{}

func (tk TransactedUniqueKeyer) GetKey(sz *sku.TransactedZettel) string {
	if sz == nil {
		return ""
	}

	return collections.MakeKey(
		sz.Kopf,
		sz.GetTai(),
		sz.TransactionIndex,
		sz.GetKennung(),
		sz.ObjekteSha,
	)
}

func MakeMutableSetUnique(c int) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.TransactedZettel](
		TransactedUniqueKeyer{},
	)
}

type TransactedHinweisKeyer struct{}

func (tk TransactedHinweisKeyer) GetKey(sz *sku.TransactedZettel) string {
	if sz == nil {
		return ""
	}

	return collections.MakeKey(
		sz.GetKennung(),
	)
}

func MakeMutableSetHinweis(c int) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.TransactedZettel](
		TransactedHinweisKeyer{},
	)
}

func ToSliceHinweisen(s MutableSet) (b []kennung.Hinweis) {
	b = make([]kennung.Hinweis, 0, s.Len())

	s.Each(
		func(z *sku.TransactedZettel) (err error) {
			b = append(b, z.GetKennung())

			return
		},
	)

	return
}
