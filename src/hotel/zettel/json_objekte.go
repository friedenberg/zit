package zettel

import (
	"encoding/json"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type JsonObjekte struct{}

func (f JsonObjekte) WriteTo(c FormatContextWrite) (n int64, err error) {
	enc := json.NewEncoder(c.Out)

	if err = enc.Encode(c.Zettel); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f JsonObjekte) ReadFrom(c *FormatContextRead) (n int64, err error) {
	return
}
