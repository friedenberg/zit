package sku

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/delta/heap"
	"github.com/friedenberg/zit/src/echo/kennung"
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