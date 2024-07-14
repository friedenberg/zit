package sku

import (
	"encoding/gob"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/heap"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

var (
	transactedKeyerObjectId interfaces.StringKeyer[*Transacted]
	TransactedSetEmpty     TransactedSet
	TransactedLessor       transactedLessor
	TransactedEqualer      transactedEqualer
)

func init() {
	transactedKeyerObjectId = &ids.ObjectIdKeyer[Transacted, *Transacted]{}
	gob.Register(transactedKeyerObjectId)

	TransactedSetEmpty = MakeTransactedSet()
	gob.Register(TransactedSetEmpty)
	gob.Register(MakeTransactedMutableSet())
}

type (
	TransactedSet        = interfaces.SetLike[*Transacted]
	TransactedMutableSet = interfaces.MutableSetLike[*Transacted]
	TransactedHeap       = heap.Heap[Transacted, *Transacted]

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

func MakeTransactedMutableSetObjectId() TransactedMutableSet {
	return collections_value.MakeMutableValueSet(
		ids.ObjectIdKeyer[Transacted, *Transacted]{},
	)
}

func MakeCheckedOutLikeMutableSet() CheckedOutLikeMutableSet {
	return collections_value.MakeMutableValueSet[CheckedOutLike](
		nil,
		// KennungKeyer[CheckedOut, *CheckedOut]{},
	)
}
