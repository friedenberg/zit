package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
)

type formatterTypActionNames struct {
	konfig               konfig_compiled.Compiled
	includeKonfigActions bool
}

func MakeFormatterTypActionNames(
	konfig konfig_compiled.Compiled,
	includeKonfigActions bool,
) *formatterTypActionNames {
	return &formatterTypActionNames{
		konfig:               konfig,
		includeKonfigActions: includeKonfigActions,
	}
}

func (e formatterTypActionNames) Format(
	w io.Writer,
	c ObjekteFormatterContext,
) (n int64, err error) {
	e1 := typ.MakeFormatterActionNames()

	ct := e.konfig.GetTyp(c.Zettel.Typ)

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

	e2 := konfig.MakeFormatterActionNames()

	if n, err = e2.Format(w, e.konfig.Actions); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *formatterTypActionNames) ReadFrom(c *ObjekteParserContext) (n int64, err error) {
	return
}
