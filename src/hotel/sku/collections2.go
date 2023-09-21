package sku

import (
	"encoding/gob"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
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
)

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
