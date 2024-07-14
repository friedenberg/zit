package zettel

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type formatterTypFormatterUTIGroups struct {
	sku.OneReader
	typBlobGetterPutter interfaces.BlobGetterPutter[*type_blobs.V0]
}

func MakeFormatterTypFormatterUTIGroups(
	sr sku.OneReader,
	tagp interfaces.BlobGetterPutter[*type_blobs.V0],
) *formatterTypFormatterUTIGroups {
	return &formatterTypFormatterUTIGroups{
		OneReader:           sr,
		typBlobGetterPutter: tagp,
	}
}

func (e formatterTypFormatterUTIGroups) Format(
	w io.Writer,
	z *sku.Transacted,
) (n int64, err error) {
	e1 := type_blobs.MakeFormatterFormatterUTIGroups()

	var skuTyp *sku.Transacted

	if skuTyp, err = e.ReadTransactedFromObjectId(z.Metadatei.GetType()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ta *type_blobs.V0

	if ta, err = e.typBlobGetterPutter.GetBlob(skuTyp.GetBlobSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer e.typBlobGetterPutter.PutBlob(ta)

	if n, err = e1.Format(w, ta); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
