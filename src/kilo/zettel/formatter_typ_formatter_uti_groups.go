package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/kilo/konfig"
)

type formatterTypFormatterUTIGroups struct {
	erworben            konfig.Compiled
	typAkteGetterPutter schnittstellen.AkteGetterPutter[*typ_akte.Akte]
}

func MakeFormatterTypFormatterUTIGroups(
	erworben konfig.Compiled,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.Akte],
) *formatterTypFormatterUTIGroups {
	return &formatterTypFormatterUTIGroups{
		erworben:            erworben,
		typAkteGetterPutter: tagp,
	}
}

func (e formatterTypFormatterUTIGroups) Format(
	w io.Writer,
	z *transacted.Zettel,
) (n int64, err error) {
	e1 := typ_akte.MakeFormatterFormatterUTIGroups()

	ct := e.erworben.GetApproximatedTyp(
		z.GetMetadatei().GetTyp(),
	).ApproximatedOrActual()

	if ct == nil {
		return
	}

	var ta *typ_akte.Akte

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
