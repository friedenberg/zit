package sku

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/delta/heap"
	"github.com/friedenberg/zit/src/echo/kennung"
)

var (
	transactedKeyerKennung schnittstellen.StringKeyerPtr[Transacted, *Transacted]
	TransactedSetEmpty     TransactedSet
	TransactedLessor       transactedLessor
)

func init() {
	transactedKeyerKennung = &KennungKeyer[Transacted, *Transacted]{}
	gob.Register(transactedKeyerKennung)

	TransactedSetEmpty = MakeTransactedSet()
	gob.Register(TransactedSetEmpty)
	gob.Register(MakeTransactedMutableSet())
}

type (
	TransactedSet        = schnittstellen.SetPtrLike[Transacted, *Transacted]
	TransactedMutableSet = schnittstellen.MutableSetPtrLike[Transacted, *Transacted]
	TransactedHeap       = heap.Heap[Transacted, *Transacted]
	CheckedOutSet        = schnittstellen.SetLike[*CheckedOut]
	CheckedOutMutableSet = schnittstellen.MutableSetLike[*CheckedOut]
)

func MakeTransactedHeap() TransactedHeap {
	return heap.Make[Transacted, *Transacted](transactedEqualer{}, transactedLessor{}, transactedResetter{})
}

func MakeTransactedSet() TransactedSet {
	return collections_ptr.MakeValueSet(transactedKeyerKennung)
}

func MakeTransactedMutableSet() TransactedMutableSet {
	return collections_ptr.MakeMutableValueSet(transactedKeyerKennung)
}

type kennungGetter interface {
	GetKennungLike() kennung.Kennung
}

type KennungKeyer[
	T kennungGetter,
	TPtr interface {
		schnittstellen.Ptr[T]
		kennungGetter
	},
] struct{}

func (sk KennungKeyer[T, TPtr]) GetKey(e T) string {
	return e.GetKennungLike().String()
}

func (sk KennungKeyer[T, TPtr]) GetKeyPtr(e TPtr) string {
	if e == nil {
		return ""
	}

	return e.GetKennungLike().String()
}
