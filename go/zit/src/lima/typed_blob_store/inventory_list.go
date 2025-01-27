package typed_blob_store

import (
	"io"
	"iter"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
)

type InventoryList struct {
	envRepo        env_repo.Env
	objectFormat   object_inventory_format.Format
	boxFormat      *box_format.BoxTransacted
	v0             inventory_list_blobs.V0
	v1             inventory_list_blobs.V1
	streamDecoders map[string]interfaces.DecoderFrom[func(*sku.Transacted) bool]
}

func MakeInventoryStore(
	dirLayout env_repo.Env,
	objectFormat object_inventory_format.Format,
	boxFormat *box_format.BoxTransacted,
) InventoryList {
	objectOptions := object_inventory_format.Options{Tai: true}

	s := InventoryList{
		envRepo:      dirLayout,
		objectFormat: objectFormat,
		boxFormat:    boxFormat,
		v0: inventory_list_blobs.MakeV0(
			objectFormat,
			objectOptions,
		),
		v1: inventory_list_blobs.V1{
			Box: boxFormat,
		},
	}

	s.streamDecoders = map[string]interfaces.DecoderFrom[func(*sku.Transacted) bool]{
		"": inventory_list_blobs.V0IterDecoder{
			V0: s.v0,
		},
		builtin_types.InventoryListTypeV1: inventory_list_blobs.V1IterDecoder{
			V1: s.v1,
		},
	}

	return s
}

func (a InventoryList) GetCommonStore() sku.BlobStore[*sku.List] {
	return a
}

func (a InventoryList) GetObjectFormat() object_inventory_format.Format {
	return a.objectFormat
}

func (a InventoryList) GetBoxFormat() *box_format.BoxTransacted {
	return a.boxFormat
}

func (a InventoryList) GetTransactedWithBlob(
	tg sku.TransactedGetter,
) (twb sku.TransactedWithBlob[*sku.List], n int64, err error) {
	twb.Transacted = tg.GetSku()
	blobSha := twb.GetBlobSha()

	var rc interfaces.ShaReadCloser

	if rc, err = a.envRepo.BlobReader(blobSha); err != nil {
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

func (a InventoryList) GetTransactedWithBlobFromReader(
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

func (a InventoryList) WriteObjectToWriter(
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

func (a InventoryList) WriteBlobToWriter(
	tipe ids.Type,
	list sku.Collection,
	w io.Writer,
) (n int64, err error) {
	switch tipe.String() {
	case "", builtin_types.InventoryListTypeV0:
		if n, err = a.v0.WriteInventoryListBlob(
			list,
			w,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case builtin_types.InventoryListTypeV1:
		if n, err = a.v1.WriteInventoryListBlob(
			list,
			w,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (a InventoryList) PutTransactedWithBlob(
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

type iterSku = func(*sku.Transacted) bool

func (a InventoryList) StreamInventoryListBlobSkus(
	tg sku.TransactedGetter,
) iter.Seq2[*sku.Transacted, error] {
	return func(yield func(*sku.Transacted, error) bool) {
		sk := tg.GetSku()
		tipe := sk.GetType()
		blobSha := sk.GetBlobSha()

		var readCloser interfaces.ShaReadCloser

		{
			var err error

			if readCloser, err = a.envRepo.BlobReader(blobSha); err != nil {
				yield(nil, errors.Wrap(err))
				return
			}
		}

		defer errors.DeferredYieldCloser(yield, readCloser)

		iter := a.IterInventoryListBlobSkusFromReader(
			tipe,
			readCloser,
		)

		for sk, err := range iter {
			if !yield(sk, err) {
				return
			}
		}
	}
}

func (a InventoryList) AllDecodedObjectsFromStream(
	reader io.Reader,
) iter.Seq2[*sku.Transacted, error] {
	return func(yield func(*sku.Transacted, error) bool) {
		decoder := triple_hyphen_io.Decoder[*ids.TypeWithObject[iterSku]]{
			Metadata: ids.TypedMetadataCoder[iterSku]{},
			Blob: ids.TypedDecodersWithoutType[iterSku](
				a.streamDecoders,
			),
		}

		if _, err := decoder.DecodeFrom(
			&ids.TypeWithObject[iterSku]{
				Object: func(sk *sku.Transacted) bool {
					return yield(sk, nil)
				},
			},
			reader,
		); err != nil {
			yield(nil, err)
			return
		}
	}
}

func (a InventoryList) IterInventoryListBlobSkusFromReader(
	tipe ids.Type,
	reader io.Reader,
) iter.Seq2[*sku.Transacted, error] {
	return func(yield func(*sku.Transacted, error) bool) {
		decoder := ids.TypedDecodersWithoutType[iterSku](
			a.streamDecoders,
		)

		if _, err := decoder.DecodeFrom(
			&ids.TypeWithObject[iterSku]{
				Type: &tipe,
				Object: func(sk *sku.Transacted) bool {
					return yield(sk, nil)
				},
			},
			reader,
		); err != nil {
			yield(nil, errors.Wrap(err))
			return
		}
	}
}

func (a InventoryList) ReadInventoryListObject(
	tipe ids.Type,
	r io.Reader,
) (out *sku.Transacted, err error) {
	switch tipe.String() {
	case "", builtin_types.InventoryListTypeV0:
		if _, out, err = a.v0.ReadInventoryListObject(
			r,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case builtin_types.InventoryListTypeV1:
		if err = a.v1.StreamInventoryListBlobSkus(
			r,
			func(sk *sku.Transacted) (err error) {
				if out == nil {
					out = sk.CloneTransacted()
				} else {
					err = errors.Errorf("expected only one sku.Transacted, but read more than one")
					return
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
