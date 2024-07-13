package zettel

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type formatterTypFormatterUTIGroups struct {
	sku.OneReader
	typAkteGetterPutter interfaces.BlobGetterPutter[*typ_akte.V0]
}

func MakeFormatterTypFormatterUTIGroups(
	sr sku.OneReader,
	tagp interfaces.BlobGetterPutter[*typ_akte.V0],
) *formatterTypFormatterUTIGroups {
	return &formatterTypFormatterUTIGroups{
		OneReader:           sr,
		typAkteGetterPutter: tagp,
	}
}

func (e formatterTypFormatterUTIGroups) Format(
	w io.Writer,
	z *sku.Transacted,
) (n int64, err error) {
	e1 := typ_akte.MakeFormatterFormatterUTIGroups()

	var skuTyp *sku.Transacted

	if skuTyp, err = e.ReadTransactedFromKennung(z.Metadatei.GetTyp()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ta *typ_akte.V0

	if ta, err = e.typAkteGetterPutter.GetBlob(skuTyp.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer e.typAkteGetterPutter.PutBlob(ta)

	if n, err = e1.Format(w, ta); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
