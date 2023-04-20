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
	z *Transacted,
) (n int64, err error) {
	e1 := typ.MakeFormatterFormatterUTIGroups()

	ct := e.erworben.GetApproximatedTyp(
		z.Objekte.Metadatei.GetTyp(),
	).ApproximatedOrActual()

	if ct == nil {
		return
	}

	if n, err = e1.Format(w, ct); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
