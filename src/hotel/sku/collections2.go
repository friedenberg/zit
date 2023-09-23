package sku

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/delta/heap"
	"github.com/friedenberg/zit/src/echo/kennung"
)

var (
	transactedKeyerKennung schnittstellen.StringKeyerPtr[Transacted2, *Transacted2]
	TransactedSetEmpty     TransactedSet
	TransactedLessor       schnittstellen.Lessor2[Transacted2, *Transacted2]
)

func init() {
	transactedKeyerKennung = &KennungKeyer[Transacted2, *Transacted2]{}
	gob.Register(transactedKeyerKennung)

	TransactedSetEmpty = MakeTransactedSet()
	gob.Register(TransactedSetEmpty)
	gob.Register(MakeTransactedMutableSet())

	TransactedLessor = Lessor[Transacted2, *Transacted2]{}
}

type (
	TransactedSet        = schnittstellen.SetPtrLike[Transacted2, *Transacted2]
	TransactedMutableSet = schnittstellen.MutableSetPtrLike[Transacted2, *Transacted2]
	TransactedHeap       = heap.Heap[Transacted2, *Transacted2]
)

func MakeTransactedHeap() TransactedHeap {
	return heap.Make[Transacted2, *Transacted2](equaler{}, lessor{}, resetter{})
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

type lessor struct{}

func (_ lessor) Less(a, b Transacted2) bool {
	return a.GetTai().Less(b.GetTai())
}

func (_ lessor) LessPtr(a, b *Transacted2) bool {
	return a.GetTai().Less(b.GetTai())
}

type equaler struct{}

func (_ equaler) Equals(a, b Transacted2) bool {
	return a.EqualsSkuLike(b)
}

func (_ equaler) EqualsPtr(a, b *Transacted2) bool {
	return a.EqualsSkuLike(b)
}

type resetter struct{}

func (_ resetter) Reset(a *Transacted2) {
	a.Reset()
}

func (_ resetter) ResetWith(a *Transacted2, b Transacted2) {
	a = &b
}

func (_ resetter) ResetWithPtr(a *Transacted2, b *Transacted2) {
	a.Kopf = b.Kopf
	a.ObjekteSha = b.ObjekteSha
	errors.PanicIfError(a.Kennung.ResetWithKennung(b.Kennung))
	a.Metadatei.ResetWith(b.Metadatei)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}
