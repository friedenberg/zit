package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/transacted"
)

type HeapTransacted = collections.Heap[transacted.Zettel, *transacted.Zettel]

func MakeHeapTransacted() HeapTransacted {
	return collections.MakeHeap[transacted.Zettel, *transacted.Zettel]()
}

type (
	MutableSet = schnittstellen.MutableSetLike[*transacted.Zettel]
)

type TransactedUniqueKeyer struct{}

func (tk TransactedUniqueKeyer) GetKey(sz *transacted.Zettel) string {
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
	return collections_value.MakeMutableValueSet[*transacted.Zettel](
		TransactedUniqueKeyer{},
	)
}

type TransactedHinweisKeyer struct{}

func (tk TransactedHinweisKeyer) GetKey(sz *transacted.Zettel) string {
	if sz == nil {
		return ""
	}

	return collections.MakeKey(
		sz.GetKennung(),
	)
}

func MakeMutableSetHinweis(c int) MutableSet {
	return collections_value.MakeMutableValueSet[*transacted.Zettel](
		TransactedHinweisKeyer{},
	)
}

func ToSliceHinweisen(s MutableSet) (b []kennung.Hinweis) {
	b = make([]kennung.Hinweis, 0, s.Len())

	s.Each(
		func(z *transacted.Zettel) (err error) {
			b = append(b, z.GetKennung())

			return
		},
	)

	return
}
