package sku

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/heap"
	"github.com/friedenberg/zit/src/echo/kennung"
)

var (
	transactedKeyerKennung schnittstellen.StringKeyerPtr[Transacted, *Transacted]
	TransactedSetEmpty     TransactedSet
	TransactedLessor       schnittstellen.Lessor2[Transacted, *Transacted]
	TransactedReseter      resetter
)

func init() {
	transactedKeyerKennung = &KennungKeyer[Transacted, *Transacted]{}
	gob.Register(transactedKeyerKennung)

	TransactedSetEmpty = MakeTransactedSet()
	gob.Register(TransactedSetEmpty)
	gob.Register(MakeTransactedMutableSet())

	TransactedLessor = &lessor{}
}

type (
	TransactedSet        = schnittstellen.SetPtrLike[Transacted, *Transacted]
	TransactedMutableSet = schnittstellen.MutableSetPtrLike[Transacted, *Transacted]
	TransactedHeap       = heap.Heap[Transacted, *Transacted]
	CheckedOutSet        = schnittstellen.SetPtrLike[CheckedOut, *CheckedOut]
	CheckedOutMutableSet = schnittstellen.MutableSetPtrLike[CheckedOut, *CheckedOut]
)

func MakeTransactedHeap() TransactedHeap {
	return heap.Make[Transacted, *Transacted](equaler{}, lessor{}, resetter{})
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

func (lessor) Less(a, b Transacted) bool {
	return a.GetTai().Less(b.GetTai())
}

func (lessor) LessPtr(a, b *Transacted) bool {
	return a.GetTai().Less(b.GetTai())
}

type equaler struct{}

func (equaler) Equals(a, b Transacted) bool {
	panic("not supported")
}

func (equaler) EqualsPtr(a, b *Transacted) bool {
	return a.EqualsSkuLikePtr(b)
}

type resetter struct{}

func (resetter) Reset(a *Transacted) {
	a.Kopf.Reset()
	a.ObjekteSha.Reset()
	a.Kennung.SetGattung(gattung.Unknown)
	a.Metadatei.Reset()
	a.TransactionIndex.Reset()
}

func (r resetter) ResetWith(a *Transacted, b Transacted) {
	a.Kopf = b.Kopf
	a.ObjekteSha = b.ObjekteSha
	a.Kennung.ResetWithKennung(b.Kennung)
	a.Metadatei.ResetWith(b.Metadatei)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}

func (resetter) ResetWithPtr(a *Transacted, b *Transacted) {
	a.Kopf = b.Kopf
	a.ObjekteSha = b.ObjekteSha
	errors.PanicIfError(a.Kennung.ResetWithKennung(&b.Kennung))
	a.Metadatei.ResetWith(b.Metadatei)
	a.TransactionIndex.SetInt(b.TransactionIndex.Int())
}
