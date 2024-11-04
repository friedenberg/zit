package blob_store

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_blobs"
)

type InventoryStore struct {
	dirLayout dir_layout.DirLayout
	v0        inventory_list_blobs.V0
	v1        inventory_list_blobs.V1
}

func MakeInventoryStore(
	dirLayout dir_layout.DirLayout,
	objectFormat object_inventory_format.Format,
	boxFormat *box_format.Box,
) InventoryStore {
	objectOptions := object_inventory_format.Options{Tai: true}

	return InventoryStore{
		dirLayout: dirLayout,
		v0: inventory_list_blobs.MakeV0(
			objectFormat,
			objectOptions,
		),
		v1: inventory_list_blobs.V1{
			Box: boxFormat,
		},
	}
}

func (a InventoryStore) GetCommonStore() CommonStore2[*sku.List] {
	return a
}

func (a InventoryStore) GetTransactedWithBlob(
	tg sku.TransactedGetter,
) (twb sku.TransactedWithBlob[*sku.List], n int64, err error) {
	sk := tg.GetSku()
	blobSha := sk.GetBlobSha()

	var rc interfaces.ShaReadCloser

	if rc, err = a.dirLayout.BlobReader(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, rc)

	if n, err = a.GetTransactedWithBlobFromReader(twb, rc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a InventoryStore) GetTransactedWithBlobFromReader(
	twb sku.TransactedWithBlob[*sku.List],
	r io.Reader,
) (n int64, err error) {
	tipe := twb.GetType()

	twb.Blob = sku.MakeList()

	switch tipe.String() {
	case "", builtin_types.InventoryListTypeV0:
		if err = inventory_list_blobs.ReadInventoryListBlob(
			a.v0,
			r,
			twb.Blob,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case builtin_types.InventoryListTypeV1:
		if err = inventory_list_blobs.ReadInventoryListBlob(
			a.v1,
			r,
			twb.Blob,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (a InventoryStore) WriteTransactedWithBlobToWriter(
  sk *sku.Transacted,
	w io.Writer,
) (n int64, err error) {
	tipe := sk.GetType()

	switch tipe.String() {
	case "", builtin_types.InventoryListTypeV0:
		if n, err = a.v0.WriteInventoryListObject(
			sk,
			w,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case builtin_types.InventoryListTypeV1:
		if n, err = a.v1.WriteInventoryListObject(
			sk,
			w,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (a InventoryStore) PutTransactedWithBlob(
	twb sku.TransactedWithBlob[*sku.List],
) (err error) {
	tipe := twb.GetType()

	switch tipe.String() {
	case "", builtin_types.InventoryListTypeV0:
	case builtin_types.InventoryListTypeV1:
	}

	sku.GetTransactedPool().Put(twb.Transacted)

	return
}
