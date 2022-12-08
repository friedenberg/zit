package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/konfig"
	"github.com/friedenberg/zit/src/golf/typ"
)

type EncoderTypActionNames struct {
	konfig konfig.Konfig
}

func (e EncoderTypActionNames) WriteTo(c FormatContextWrite) (n int64, err error) {
	e1 := typ.MakeFormatterActionNames(c.Out)

	//TODO-P2 move to store_objekten
	ct := e.konfig.Transacted.Objekte.GetTyp(c.Zettel.Typ.String())

	if ct == nil {
		return
	}

	ty := &typ.Transacted{
		Sku:     ct.Sku,
		Objekte: ct.Typ,
	}

	if n, err = e1.Encode(ty); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *EncoderTypActionNames) ReadFrom(c *FormatContextRead) (n int64, err error) {
	return
}
