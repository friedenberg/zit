package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/hotel/konfig_compiled"
	"github.com/friedenberg/zit/src/india/typ"
)

type EncoderTypActionNames struct {
	konfig konfig_compiled.Compiled
}

func (e EncoderTypActionNames) WriteTo(c FormatContextWrite) (n int64, err error) {
	e1 := typ.MakeFormatterActionNames(c.Out)

	//TODO-P2 move to store_objekten
	ct := e.konfig.GetTyp(c.Zettel.Typ.String())

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
