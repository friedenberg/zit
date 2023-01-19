package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/india/konfig"
)

type formatterTypActionNames struct {
	erworben             konfig.Compiled
	includeKonfigActions bool
}

func MakeFormatterTypActionNames(
	erworben konfig.Compiled,
	includeKonfigActions bool,
) *formatterTypActionNames {
	return &formatterTypActionNames{
		erworben:             erworben,
		includeKonfigActions: includeKonfigActions,
	}
}

func (e formatterTypActionNames) Format(
	w io.Writer,
	c ObjekteFormatterContext,
) (n int64, err error) {
	e1 := typ.MakeFormatterActionNames()

	ct := e.erworben.GetTyp(c.Zettel.Typ)

	if ct == nil {
		return
	}

	if n, err = e1.Format(w, ct); err != nil {
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

func (f *formatterTypActionNames) ReadFrom(c *ObjekteParserContext) (n int64, err error) {
	return
}
