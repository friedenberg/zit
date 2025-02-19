package inventory_list_store

import (
	"io"
	"iter"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
)

// TODO add triple_hyphen_io coder
type objectBlobStoreV1 struct {
	pathLog        string
	blobType       ids.Type
	typedBlobStore typed_blob_store.InventoryList
	blobStore      interfaces.LocalBlobStore
}

func (store *objectBlobStoreV1) getType() ids.Type {
	return store.blobType
}

func (store *objectBlobStoreV1) getTypedBlobStore() typed_blob_store.InventoryList {
	return store.typedBlobStore
}

func (store *objectBlobStoreV1) ReadOneSha(
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

func (store *objectBlobStoreV1) WriteInventoryListObject(
	object *sku.Transacted,
) (err error) {
	var blobStoreWriteCloser interfaces.ShaWriteCloser

	if blobStoreWriteCloser, err = store.blobStore.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, blobStoreWriteCloser)

	object.Metadata.Type = store.blobType

	var file *os.File

	if file, err = files.OpenExclusiveWriteOnlyAppend(store.pathLog); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, file)
	defer errors.Deferred(&err, file.Sync)

	if _, err = store.typedBlobStore.WriteObjectToWriter(
		store.blobType,
		object,
		io.MultiWriter(blobStoreWriteCloser, file),
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

func (store *objectBlobStoreV1) IterAllInventoryLists() iter.Seq2[*sku.Transacted, error] {
	return func(yield func(*sku.Transacted, error) bool) {
		var file *os.File

		{
			var err error

			if file, err = files.OpenReadOnly(store.pathLog); err != nil {
				yield(nil, err)
				return
			}
		}

		seq := store.typedBlobStore.AllDecodedObjectsFromStream(
			file,
		)

		for sk, err := range seq {
			if err != nil {
				if !yield(nil, err) {
					return
				}
			}

			if !yield(sk, nil) {
				return
			}
		}

		if err := file.Close(); err != nil {
			yield(nil, err)
			return
		}
	}
}
