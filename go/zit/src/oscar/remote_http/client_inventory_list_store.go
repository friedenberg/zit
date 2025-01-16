package remote_http

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (client Client) WriteInventoryListObject(t *sku.Transacted) (err error) {
	return todo.Implement()
}

func (client Client) ImportInventoryList(bs interfaces.BlobStore, t *sku.Transacted) (err error) {
	return todo.Implement()
}

func (client Client) ReadLast() (max *sku.Transacted, err error) {
	return nil, todo.Implement()
}

func (client Client) StreamInventoryList(
	blobSha interfaces.Sha,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return todo.Implement()
}

func (client Client) ReadAllSkus(
	f func(besty, sk *sku.Transacted) error,
) (err error) {
	return todo.Implement()
}

func (client Client) ReadAllInventoryLists(
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return todo.Implement()
}
