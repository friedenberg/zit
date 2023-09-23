package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/delta/heap"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type HeapTransacted = heap.Heap[sku.Transacted2, *sku.Transacted2]

func MakeHeapTransacted() HeapTransacted {
	return heap.Make[sku.Transacted2, *sku.Transacted2](
		sku.Equaler[sku.Transacted2, *sku.Transacted2]{},
		sku.Lessor[sku.Transacted2, *sku.Transacted2]{},
		sku.Resetter[sku.Transacted2, *sku.Transacted2]{},
	)
}

func MakeHeapTransactedReversed() HeapTransacted {
	return heap.Make[sku.Transacted2, *sku.Transacted2](
		sku.Equaler[sku.Transacted2, *sku.Transacted2]{},
		values.ReverseLessor[sku.Transacted2, *sku.Transacted2]{
			Inner: sku.Lessor[sku.Transacted2, *sku.Transacted2]{},
		},
		sku.Resetter[sku.Transacted2, *sku.Transacted2]{},
	)
}

type (
	MutableSet = schnittstellen.MutableSetLike[*sku.Transacted2]
)

type TransactedUniqueKeyer struct{}

func (tk TransactedUniqueKeyer) GetKey(sz *sku.Transacted2) string {
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
	return collections_value.MakeMutableValueSet[*sku.Transacted2](
		TransactedUniqueKeyer{},
	)
}

type TransactedHinweisKeyer struct{}

func (tk TransactedHinweisKeyer) GetKey(sz *sku.Transacted2) string {
	if sz == nil {
		return ""
	}

	return collections.MakeKey(
		sz.GetKennung(),
	)
}

func MakeMutableSetHinweis(c int) MutableSet {
	return collections_value.MakeMutableValueSet[*sku.Transacted2](
		TransactedHinweisKeyer{},
	)
}
