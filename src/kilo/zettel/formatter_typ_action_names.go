package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/kilo/konfig"
)

type formatterTypActionNames struct {
	erworben             konfig.Compiled
	includeKonfigActions bool
	typAkteGetterPutter  schnittstellen.AkteGetterPutter[*typ_akte.V0]
}

func MakeFormatterTypActionNames(
	erworben konfig.Compiled,
	includeKonfigActions bool,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) *formatterTypActionNames {
	return &formatterTypActionNames{
		erworben:             erworben,
		includeKonfigActions: includeKonfigActions,
		typAkteGetterPutter:  tagp,
	}
}

func (e formatterTypActionNames) Format(
	w io.Writer,
	z *sku.Transacted,
) (n int64, err error) {
	e1 := typ_akte.MakeFormatterActionNames()

	ct := e.erworben.GetApproximatedTyp(z.GetMetadatei().GetTyp())
	t := ct.ApproximatedOrActual()

	if t == nil {
		return
	}

	var ta *typ_akte.V0

	if ta, err = e.typAkteGetterPutter.GetAkte(t.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer e.typAkteGetterPutter.PutAkte(ta)

	if n, err = e1.Format(w, ta); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !e.includeKonfigActions {
		return
	}

	e2 := erworben.MakeFormatterActionNames()

	if n, err = e2.Format(w, e.erworben.Actions); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
