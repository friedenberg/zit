package inventory_list_store

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
)

type objectBlobStoreV0 struct {
	lock           sync.Mutex
	blobType       ids.Type
	typedBlobStore typed_blob_store.InventoryList
	blobStore      interfaces.LocalBlobStore
}

func (store *objectBlobStoreV0) getType() ids.Type {
	return store.blobType
}

func (store *objectBlobStoreV0) getTypedBlobStore() typed_blob_store.InventoryList {
	return store.typedBlobStore
}

func (store *objectBlobStoreV0) ReadOneSha(
	id interfaces.Stringer,
) (object *sku.Transacted, err error) {
	var sh sha.Sha

	if err = sh.Set(id.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var readCloser sha.ReadCloser

	if readCloser, err = store.blobStore.BlobReader(&sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, readCloser)

	if object, err = store.typedBlobStore.ReadInventoryListObject(
		store.blobType,
		readCloser,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (store *objectBlobStoreV0) WriteInventoryListObject(
	object *sku.Transacted,
) (err error) {
	store.lock.Lock()
	defer store.lock.Unlock()

	var blobStoreWriteCloser interfaces.ShaWriteCloser

	if blobStoreWriteCloser, err = store.blobStore.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, blobStoreWriteCloser)

	object.Metadata.Type = store.blobType

	if _, err = store.typedBlobStore.WriteObjectToWriter(
		object,
		blobStoreWriteCloser,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = object.CalculateObjectShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ui.Log().Printf(
		"saved inventory list: %q",
		sku.String(object),
	)

	return
}
