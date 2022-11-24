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
	e1 := typ.MakeEncoderActionNames(c.Out, e.konfig)

	if n, err = e1.Encode(&c.Zettel.Typ); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *EncoderTypActionNames) ReadFrom(c *FormatContextRead) (n int64, err error) {
	return
}
