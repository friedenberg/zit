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
	dirLayout  dir_layout.DirLayout
	boxFormat  *box_format.Box
	v0         inventory_list_blobs.V0
	v1         inventory_list_blobs.V1
	listFormat sku.ListFormat
}

func MakeInventoryStore(
	dirLayout dir_layout.DirLayout,
	objectFormat object_inventory_format.Format,
	boxFormat *box_format.Box,
) InventoryStore {
	objectOptions := object_inventory_format.Options{Tai: true}

	s := InventoryStore{
		dirLayout: dirLayout,
		boxFormat: boxFormat,
		v0: inventory_list_blobs.MakeV0(
			objectFormat,
			objectOptions,
		),
		v1: inventory_list_blobs.V1{
			Box: boxFormat,
		},
	}

	s.listFormat = s.formatForVersion(objectFormat, objectOptions)

	return s
}

func (s *InventoryStore) formatForVersion(
	objectFormat object_inventory_format.Format,
	objectOptions object_inventory_format.Options,
) sku.ListFormat {
	v := s.dirLayout.GetStoreVersion().GetInt()

	switch {
	case v <= 6:
		return inventory_list_blobs.MakeV0(
			objectFormat,
			objectOptions,
		)

	default:
		return inventory_list_blobs.V1{
			Box: s.boxFormat,
		}
	}
}

func (a InventoryStore) GetCommonStore() CommonStore2[*sku.List] {
	return a
}

func (a InventoryStore) GetListFormat() sku.ListFormat {
	return a.listFormat
}

func (a InventoryStore) GetTransactedWithBlob(
	tg sku.TransactedGetter,
) (twb sku.TransactedWithBlob[*sku.List], n int64, err error) {
	twb.Transacted = tg.GetSku()
	blobSha := twb.GetBlobSha()

	var rc interfaces.ShaReadCloser

	if rc, err = a.dirLayout.BlobReader(blobSha); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, rc)

	if n, err = a.GetTransactedWithBlobFromReader(&twb, rc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (a InventoryStore) GetTransactedWithBlobFromReader(
	twb *sku.TransactedWithBlob[*sku.List],
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

func (a InventoryStore) WriteObjectToWriter(
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

func (a InventoryStore) WriteBlobToWriter(
	sk *sku.Transacted,
  b *sku.List,
	w io.Writer,
) (n int64, err error) {
	tipe := sk.GetType()

	switch tipe.String() {
	case "", builtin_types.InventoryListTypeV0:
		if n, err = a.v0.WriteInventoryListBlob(
			b,
			w,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case builtin_types.InventoryListTypeV1:
		if n, err = a.v1.WriteInventoryListBlob(
			b,
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
