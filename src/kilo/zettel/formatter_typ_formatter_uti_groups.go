package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
)

type formatterTypFormatterUTIGroups struct {
	konfig konfig_compiled.Compiled
}

func MakeFormatterTypFormatterUTIGroups(
	konfig konfig_compiled.Compiled,
) *formatterTypFormatterUTIGroups {
	return &formatterTypFormatterUTIGroups{
		konfig: konfig,
	}
}

func (e formatterTypFormatterUTIGroups) Format(
	w io.Writer,
	c ObjekteFormatterContext,
) (n int64, err error) {
	e1 := typ.MakeFormatterFormatterUTIGroups()

	ct := e.konfig.GetTyp(c.Zettel.Typ)

	if ct == nil {
		return
	}

	if n, err = e1.Format(w, ct); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *formatterTypFormatterUTIGroups) ReadFrom(c *ObjekteParserContext) (n int64, err error) {
	return
}
