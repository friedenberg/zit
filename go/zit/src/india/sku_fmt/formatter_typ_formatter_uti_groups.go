package sku_fmt

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/blob_store"
)

type formatterTypFormatterUTIGroups struct {
	sku.OneReader
	blob_store.TypeStore
}

func MakeFormatterTypFormatterUTIGroups(
	sr sku.OneReader,
	typeBlobStore blob_store.TypeStore,
) *formatterTypFormatterUTIGroups {
	return &formatterTypFormatterUTIGroups{
		OneReader: sr,
		TypeStore: typeBlobStore,
	}
}

func (e formatterTypFormatterUTIGroups) Format(
	w io.Writer,
	z *sku.Transacted,
) (n int64, err error) {
	var skuTyp *sku.Transacted

	if skuTyp, err = e.ReadTransactedFromObjectId(z.Metadata.GetType()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ta type_blobs.Blob

	if ta, _, err = e.ParseTypedBlob(
		skuTyp.GetType(),
		skuTyp.GetBlobSha(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer e.PutTypedBlob(skuTyp.GetType(), ta)

	for groupName, group := range ta.GetFormatterUTIGroups() {
		sb := bytes.NewBuffer(nil)

		sb.WriteString(groupName)

		for uti, formatter := range group.Map() {
			sb.WriteString(" ")
			sb.WriteString(uti)
			sb.WriteString(" ")
			sb.WriteString(formatter)
		}

		sb.WriteString("\n")

		if n, err = io.Copy(w, sb); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
