package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
)

type EncoderTypActionNames struct {
	konfig konfig_compiled.Compiled
}

func (e EncoderTypActionNames) Format(
	w io.Writer,
	c FormatContextWrite,
) (n int64, err error) {
	e1 := typ.MakeFormatterActionNames(w)

	ct := e.konfig.GetTyp(c.Zettel.Typ)

	if ct == nil {
		return
	}

	if n, err = e1.Encode(ct); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *EncoderTypActionNames) ReadFrom(c *FormatContextRead) (n int64, err error) {
	return
}
