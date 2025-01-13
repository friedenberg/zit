package sku

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/heap"
)

type InventoryListStore interface {
	WriteInventoryList(t InventoryList) (err error)
	// ReadInventoryList(ids.Tai) (*sku.Transacted, *sku.List, error)

	// ReadAllSkus(
	// 	f func(besty, sk *sku.Transacted) error,
	// ) (err error)

	// ReadAllInventoryListsSince(
	// since ids.Tai,
	// 	f interfaces.FuncIter[*sku.Transacted],
	// ) (err error)

	// ReadAllInventoryLists(
	// 	f interfaces.FuncIter[*sku.Transacted],
	// ) (err error)
}

type ListFormat interface {
	GetListFormat() ListFormat
	WriteInventoryListBlob(Collection, io.Writer) (int64, error)
	WriteInventoryListObject(*Transacted, io.Writer) (int64, error)
	ReadInventoryListObject(io.Reader) (int64, *Transacted, error)
	StreamInventoryListBlobSkus(
		rf io.Reader,
		f interfaces.FuncIter[*Transacted],
	) error
}

type List = heap.Heap[Transacted, *Transacted]

func MakeList() *List {
	h := heap.Make(
		transactedEqualer{},
		transactedLessor{},
		transactedResetter{},
	)

	h.SetPool(GetTransactedPool())

	return h
}

var ResetterList resetterList

type resetterList struct{}

func (resetterList) Reset(a *List) {
	a.Reset()
}

func (resetterList) ResetWith(a, b *List) {
	a.ResetWith(b)
}
