package sku_fmt

import (
	"bytes"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
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
	var skuTyp *sku.Transacted

	if skuTyp, err = e.ReadTransactedFromObjectId(z.Metadata.GetType()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ta *type_blobs.V0

	if ta, err = e.typBlobGetterPutter.GetBlob(skuTyp.GetBlobSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer e.typBlobGetterPutter.PutBlob(ta)

	for groupName, group := range ta.FormatterUTIGroups {
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
