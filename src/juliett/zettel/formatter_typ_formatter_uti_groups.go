package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/transacted"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/india/konfig"
)

type formatterTypFormatterUTIGroups struct {
	erworben            konfig.Compiled
	typAkteGetterPutter schnittstellen.AkteGetterPutter[*typ.Akte]
}

func MakeFormatterTypFormatterUTIGroups(
	erworben konfig.Compiled,
	tagp schnittstellen.AkteGetterPutter[*typ.Akte],
) *formatterTypFormatterUTIGroups {
	return &formatterTypFormatterUTIGroups{
		erworben:            erworben,
		typAkteGetterPutter: tagp,
	}
}

func (e formatterTypFormatterUTIGroups) Format(
	w io.Writer,
	z *transacted.Zettel) (n int64, err error) {
	e1 := typ.MakeFormatterFormatterUTIGroups()

	ct := e.erworben.GetApproximatedTyp(
		z.GetMetadatei().GetTyp(),
	).ApproximatedOrActual()

	if ct == nil {
		return
	}

	var ta *typ.Akte

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
