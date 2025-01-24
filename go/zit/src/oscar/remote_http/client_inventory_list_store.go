package remote_http

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (client client) FormatForVersion(
	sv interfaces.StoreVersion,
) sku.ListFormat {
	return client.localInventoryListStore.FormatForVersion(sv)
}

func (client client) WriteInventoryListObject(t *sku.Transacted) (err error) {
	return todo.Implement()
}

func (client client) ImportInventoryList(
	bs interfaces.BlobStore,
	t *sku.Transacted,
) (err error) {
	return todo.Implement()
}

func (client client) ReadLast() (max *sku.Transacted, err error) {
	return nil, todo.Implement()
}

func (client client) StreamInventoryList(
	blobSha interfaces.Sha,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return todo.Implement()
}

func (client client) ReadAllSkus(
	f func(besty, sk *sku.Transacted) error,
) (err error) {
	return todo.Implement()
}

func (client client) ReadAllInventoryLists(
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return todo.Implement()
}
