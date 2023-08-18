package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/india/konfig"
)

type formatterTypActionNames struct {
	erworben             konfig.Compiled
	includeKonfigActions bool
	typAkteGetterPutter  schnittstellen.AkteGetterPutter[*typ.Akte]
}

func MakeFormatterTypActionNames(
	erworben konfig.Compiled,
	includeKonfigActions bool,
	tagp schnittstellen.AkteGetterPutter[*typ.Akte],
) *formatterTypActionNames {
	return &formatterTypActionNames{
		erworben:             erworben,
		includeKonfigActions: includeKonfigActions,
		typAkteGetterPutter:  tagp,
	}
}

func (e formatterTypActionNames) Format(
	w io.Writer,
	z *Transacted,
) (n int64, err error) {
	e1 := typ.MakeFormatterActionNames()

	ct := e.erworben.GetApproximatedTyp(z.GetMetadatei().GetTyp())
	t := ct.ApproximatedOrActual()

	if t == nil {
		return
	}

	var ta *typ.Akte

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
