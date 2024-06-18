package zettel

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type formatterTypFormatterUTIGroups struct {
	sku.OneReader
	typAkteGetterPutter schnittstellen.AkteGetterPutter[*typ_akte.V0]
}

func MakeFormatterTypFormatterUTIGroups(
	sr sku.OneReader,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
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

	if skuTyp, err = e.ReadOne(z.Metadatei.GetTyp()); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ta *typ_akte.V0

	if ta, err = e.typAkteGetterPutter.GetAkte(skuTyp.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer e.typAkteGetterPutter.PutAkte(ta)

	if n, err = e1.Format(w, ta); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
