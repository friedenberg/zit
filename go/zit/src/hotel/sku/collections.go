package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/heap"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

var (
	transactedKeyerObjectId   ObjectIdKeyer[*Transacted]
	externalLikeKeyerObjectId = interfaces.CompoundKeyer[ExternalLike]{
		ObjectIdKeyer[ExternalLike]{},
		ExternalObjectIdKeyer[ExternalLike]{},
		DescriptionKeyer[ExternalLike]{},
	}
	checkedOutKeyerObjectId ObjectIdKeyer[*CheckedOut]
	TransactedSetEmpty      TransactedSet
	TransactedLessor        transactedLessor
	TransactedEqualer       transactedEqualer
)

func GetExternalLikeKeyer[
	T interface {
		ExternalObjectIdGetter
		ids.ObjectIdGetter
		ExternalLikeGetter
	},
]() interfaces.StringKeyer[T] {
	return interfaces.CompoundKeyer[T]{
		ObjectIdKeyer[T]{},
		ExternalObjectIdKeyer[T]{},
		DescriptionKeyer[T]{},
	}
}

type Collection interfaces.Collection[*Transacted]

func init() {
	gob.Register(transactedKeyerObjectId)

	TransactedSetEmpty = MakeTransactedSet()
	gob.Register(TransactedSetEmpty)
	gob.Register(MakeTransactedMutableSet())
}

type (
	TransactedSet        = interfaces.SetLike[*Transacted]
	TransactedMutableSet = interfaces.MutableSetLike[*Transacted]
	TransactedHeap       = heap.Heap[Transacted, *Transacted]

	ExternalLikeSet        = interfaces.SetLike[ExternalLike]
	ExternalLikeMutableSet = interfaces.MutableSetLike[ExternalLike]

	CheckedOutSet        = interfaces.SetLike[*CheckedOut]
	CheckedOutMutableSet = interfaces.MutableSetLike[*CheckedOut]
)

func MakeTransactedHeap() *TransactedHeap {
	h := heap.Make(
		transactedEqualer{},
		transactedLessor{},
		transactedResetter{},
	)

	h.SetPool(GetTransactedPool())

	return h
}

func MakeTransactedSet() TransactedSet {
	return collections_value.MakeValueSet(transactedKeyerObjectId)
}

func MakeTransactedMutableSet() TransactedMutableSet {
	return collections_value.MakeMutableValueSet(transactedKeyerObjectId)
}

func MakeExternalLikeSet() ExternalLikeSet {
	return collections_value.MakeValueSet(externalLikeKeyerObjectId)
}

func MakeExternalLikeMutableSet() ExternalLikeMutableSet {
	return collections_value.MakeMutableValueSet(externalLikeKeyerObjectId)
}

func MakeCheckedOutSet() CheckedOutSet {
	return collections_value.MakeValueSet(checkedOutKeyerObjectId)
}

func MakeCheckedOutMutableSet() CheckedOutMutableSet {
	return collections_value.MakeMutableValueSet(checkedOutKeyerObjectId)
}
