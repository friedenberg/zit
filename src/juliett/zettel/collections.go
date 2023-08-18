package zettel

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type HeapTransacted = collections.Heap[Transacted, *Transacted]

func MakeHeapTransacted() HeapTransacted {
	return collections.MakeHeap[Transacted, *Transacted]()
}

type (
	MutableSet = schnittstellen.MutableSetLike[*Transacted]
)

type TransactedUniqueKeyer struct{}

func (tk TransactedUniqueKeyer) GetKey(sz *Transacted) string {
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
	return collections_value.MakeMutableValueSet[*Transacted](
		TransactedUniqueKeyer{},
	)
}

type TransactedHinweisKeyer struct{}

func (tk TransactedHinweisKeyer) GetKey(sz *Transacted) string {
	if sz == nil {
		return ""
	}

	return collections.MakeKey(
		sz.GetKennung(),
	)
}

func MakeMutableSetHinweis(c int) MutableSet {
	return collections_value.MakeMutableValueSet[*Transacted](
		TransactedHinweisKeyer{},
	)
}

func ToSliceHinweisen(s MutableSet) (b []kennung.Hinweis) {
	b = make([]kennung.Hinweis, 0, s.Len())

	s.Each(
		func(z *Transacted) (err error) {
			b = append(b, z.GetKennung())

			return
		},
	)

	return
}
