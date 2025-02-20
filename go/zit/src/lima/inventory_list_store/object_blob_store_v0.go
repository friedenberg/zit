package inventory_list_store

import (
	"iter"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
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
		store.blobType,
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

func (s *objectBlobStoreV0) readOnePath(p string) (o *sku.Transacted, err error) {
	var sh *sha.Sha

	if sh, err = sha.MakeShaFromPath(p); err != nil {
		err = errors.Wrapf(err, "Path: %q", p)
		return
	}

	if o, err = s.ReadOneSha(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = o.CalculateObjectShas(); err != nil {
		if errors.Is(err, object_inventory_format.ErrEmptyTai) {
			var t ids.Tai
			err1 := t.Set(o.ObjectId.String())

			if err1 != nil {
				err = errors.Wrapf(err, "%s", sku.StringTaiGenreObjectIdShaBlob(o))
				return
			}

			o.SetTai(t)

			if err = o.CalculateObjectShas(); err != nil {
				err = errors.Wrapf(err, "%#v", o)
				return
			}
		} else {
			err = errors.Wrapf(err, "%#v", o)
		}

		return
	}

	return
}

func (s *objectBlobStoreV0) IterAllInventoryLists() iter.Seq2[*sku.Transacted, error] {
	return func(yield func(*sku.Transacted, error) bool) {
		for sh, err := range s.blobStore.AllBlobs() {
			if err != nil {
				if !yield(nil, err) {
					return
				}
			}

			var decodedList *sku.Transacted

			if decodedList, err = s.ReadOneSha(sh); err != nil {
				if !yield(nil, errors.Wrapf(err, "Sha: %q", sh)) {
					return
				}
			}

			if !yield(decodedList, nil) {
				return
			}
		}
	}
}
