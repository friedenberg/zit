package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/delta/heap"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

var (
	transactedKeyerKennung schnittstellen.StringKeyer[*Transacted]
	TransactedSetEmpty     TransactedSet
	TransactedLessor       transactedLessor
	TransactedEqualer      transactedEqualer
)

func init() {
	transactedKeyerKennung = &KennungKeyer[Transacted, *Transacted]{}
	gob.Register(transactedKeyerKennung)

	TransactedSetEmpty = MakeTransactedSet()
	gob.Register(TransactedSetEmpty)
	gob.Register(MakeTransactedMutableSet())
}

type (
	TransactedSet        = schnittstellen.SetLike[*Transacted]
	TransactedMutableSet = schnittstellen.MutableSetLike[*Transacted]
	TransactedHeap       = heap.Heap[Transacted, *Transacted]
	CheckedOutSet        = schnittstellen.SetLike[*CheckedOut]
	CheckedOutMutableSet = schnittstellen.MutableSetLike[*CheckedOut]
)

func MakeTransactedHeap() *TransactedHeap {
	h := heap.Make[Transacted, *Transacted](
		transactedEqualer{},
		transactedLessor{},
		transactedResetter{},
	)

	h.SetPool(GetTransactedPool())

	return h
}

func MakeTransactedSet() TransactedSet {
	return collections_value.MakeValueSet(transactedKeyerKennung)
}

func MakeTransactedMutableSet() TransactedMutableSet {
	return collections_value.MakeMutableValueSet(transactedKeyerKennung)
}

func MakeTransactedMutableSetKennung() TransactedMutableSet {
	return collections_value.MakeMutableValueSet(
		KennungKeyer[Transacted, *Transacted]{},
	)
}

type kennungGetter interface {
	GetKennung() kennung.Kennung
}

type KennungKeyer[
	T any,
	TPtr interface {
		schnittstellen.Ptr[T]
		kennungGetter
	},
] struct{}

func (sk KennungKeyer[T, TPtr]) GetKey(e TPtr) string {
	return e.GetKennung().String()
}
