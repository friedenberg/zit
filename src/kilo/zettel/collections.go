package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/delta/heap"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type HeapTransacted = heap.Heap[sku.Transacted, *sku.Transacted]

func MakeHeapTransacted() HeapTransacted {
	return heap.Make[sku.Transacted, *sku.Transacted](
		sku.Equaler[sku.Transacted, *sku.Transacted]{},
		sku.Lessor[sku.Transacted, *sku.Transacted]{},
		sku.Resetter[sku.Transacted, *sku.Transacted]{},
	)
}

func MakeHeapTransactedReversed() HeapTransacted {
	return heap.Make[sku.Transacted, *sku.Transacted](
		sku.Equaler[sku.Transacted, *sku.Transacted]{},
		values.ReverseLessor[sku.Transacted, *sku.Transacted]{
			Inner: sku.Lessor[sku.Transacted, *sku.Transacted]{},
		},
		sku.Resetter[sku.Transacted, *sku.Transacted]{},
	)
}

type (
	MutableSet = schnittstellen.MutableSetLike[*sku.Transacted]
)

type TransactedUniqueKeyer struct{}

func (tk TransactedUniqueKeyer) GetKey(sz *sku.Transacted) string {
	if sz == nil {
		return ""
	}

	return collections.MakeKey(
		sz.GetTai(),
		sz.TransactionIndex,
		sz.GetKennung(),
		sz.ObjekteSha,
	)
}

func MakeMutableSetUnique(c int) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.Transacted](
		TransactedUniqueKeyer{},
	)
}

type TransactedHinweisKeyer struct{}

func (tk TransactedHinweisKeyer) GetKey(sz *sku.Transacted) string {
	if sz == nil {
		return ""
	}

	return collections.MakeKey(
		sz.GetKennung(),
	)
}

func MakeMutableSetHinweis(c int) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.Transacted](
		TransactedHinweisKeyer{},
	)
}
