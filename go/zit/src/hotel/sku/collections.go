package sku

import (
	"encoding/gob"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/heap"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

var (
	transactedKeyerObjectId   ids.ObjectIdKeyer[*Transacted]
	externalLikeKeyerObjectId ExternalObjectIdKeyer[ExternalLike]
	checkedOutKeyerObjectId   ids.ObjectIdKeyer[*CheckedOut]
	TransactedSetEmpty        TransactedSet
	TransactedLessor          transactedLessor
	TransactedEqualer         transactedEqualer
)

type Collection interfaces.Collection[*Transacted]

type ExternalObjectIdKeyer[
	T interface {
		ids.ObjectIdGetter
		ExternalObjectIdGetter
		TransactedGetter
	},
] struct{}

func (ExternalObjectIdKeyer[T]) GetKey(el T) string {
	eoid := el.GetExternalObjectId()

	if !eoid.IsEmpty() {
		return eoid.String()
	}

	if !el.GetSku().ObjectId.IsEmpty() {
		return el.GetSku().ObjectId.String()
	}

	desc := el.GetSku().Metadata.Description.String()

	if desc != "" {
		return desc
	}

	panic(fmt.Sprintf("empty key for external like: %#v", el))
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
