package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/india/konfig"
)

type formatterTypFormatterUTIGroups struct {
	erworben konfig.Compiled
}

func MakeFormatterTypFormatterUTIGroups(
	erworben konfig.Compiled,
) *formatterTypFormatterUTIGroups {
	return &formatterTypFormatterUTIGroups{
		erworben: erworben,
	}
}

func (e formatterTypFormatterUTIGroups) Format(
	w io.Writer,
	c ObjekteFormatterContext,
) (n int64, err error) {
	e1 := typ.MakeFormatterFormatterUTIGroups()

	ct := e.erworben.GetApproximatedTyp(c.Zettel.Typ)

	if ct == nil {
		return
	}

	if n, err = e1.Format(w, ct.Unwrap()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *formatterTypFormatterUTIGroups) ReadFrom(c *ObjekteParserContext) (n int64, err error) {
	return
}
