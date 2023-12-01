package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/konfig"
)

type formatterTypFormatterUTIGroups struct {
	erworben            konfig.Compiled
	typAkteGetterPutter schnittstellen.AkteGetterPutter[*typ_akte.V0]
}

func MakeFormatterTypFormatterUTIGroups(
	erworben konfig.Compiled,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) *formatterTypFormatterUTIGroups {
	return &formatterTypFormatterUTIGroups{
		erworben:            erworben,
		typAkteGetterPutter: tagp,
	}
}

func (e formatterTypFormatterUTIGroups) Format(
	w io.Writer,
	z *sku.Transacted,
) (n int64, err error) {
	e1 := typ_akte.MakeFormatterFormatterUTIGroups()

	ct := e.erworben.GetApproximatedTyp(
		z.Metadatei.GetTyp(),
	).ApproximatedOrActual()

	if ct == nil {
		return
	}

	var ta *typ_akte.V0

	if ta, err = e.typAkteGetterPutter.GetAkte(ct.GetAkteSha()); err != nil {
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
