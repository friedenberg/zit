package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/heap"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

var (
	transactedKeyerObjectId   ids.ObjectIdKeyer[*Transacted]
	externalLikeKeyerObjectId ExternalObjectIdKeyer[ExternalLike]
	TransactedSetEmpty        TransactedSet
	TransactedLessor          transactedLessor
	TransactedEqualer         transactedEqualer
)

type ExternalObjectIdKeyer[
	T ExternalObjectIdGetter,
] struct{}

func (ExternalObjectIdKeyer[T]) GetKey(e T) string {
	return e.GetExternalObjectId().String()
}

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

	CheckedOutSet            = interfaces.SetLike[*CheckedOut]
	CheckedOutMutableSet     = interfaces.MutableSetLike[*CheckedOut]
	CheckedOutLikeSet        = interfaces.SetLike[CheckedOutLike]
	CheckedOutLikeMutableSet = interfaces.MutableSetLike[CheckedOutLike]
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

func MakeCheckedOutLikeMutableSet() CheckedOutLikeMutableSet {
	return collections_value.MakeMutableValueSet[CheckedOutLike](
		nil,
	)
}
