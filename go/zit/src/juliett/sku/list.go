package sku

import (
	"io"
	"iter"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/heap"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type InventoryListStore interface {
	FormatForVersion(sv interfaces.StoreVersion) ListFormat
	WriteInventoryListObject(t *Transacted) (err error)
	ImportInventoryList(bs interfaces.BlobStore, t *Transacted) (err error)
	// WriteInventoryListStream(list *Transacted, ) (err error)
	// ReadInventoryList(ids.Tai) (*sku.Transacted, *sku.List, error)

	ReadLast() (max *Transacted, err error)

	IterInventoryList(interfaces.Sha) iter.Seq2[*Transacted, error]

	ReadAllSkus(
		f func(besty, sk *Transacted) error,
	) (err error)

	// ReadAllInventoryListsSince(
	// since ids.Tai,
	// 	f interfaces.FuncIter[*sku.Transacted],
	// ) (err error)

	IterAllInventoryLists() iter.Seq2[*Transacted, error]
	MakeImporter(ImporterOptions, StoreOptions) Importer
	ImportList(*List, Importer) error
}

type ListFormat interface {
	GetType() ids.Type
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
		transactedLessorStable{},
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

func CollectList(seq iter.Seq2[*Transacted, error]) (list *List, err error) {
	list = MakeList()

	for sk, iterErr := range seq {
		if iterErr != nil {
			err = errors.Wrap(iterErr)
			return
		}

		list.Add(sk)
	}

	return
}
