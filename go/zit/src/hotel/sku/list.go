package sku

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

// TODO modify this to accept a interfaces.Collection[*sku.Transacted]
type ListFormat interface {
	GetListFormat() ListFormat
	WriteInventoryListBlob(*List, io.Writer) (int64, error)
	WriteInventoryListObject(*Transacted, io.Writer) (int64, error)
	ReadInventoryListObject(io.Reader) (int64, *Transacted, error)
	StreamInventoryListBlobSkus(
		rf io.Reader,
		f interfaces.FuncIter[*Transacted],
	) error
}

type List = TransactedHeap

func MakeList() *List {
	return MakeTransactedHeap()
}

var ResetterList resetterList

type resetterList struct{}

func (resetterList) Reset(a *List) {
	a.Reset()
}

func (resetterList) ResetWith(a, b *List) {
	a.ResetWith(b)
}
