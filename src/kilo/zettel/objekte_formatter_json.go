package zettel

import (
	"encoding/json"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type objekteFormatterJson struct{}

func MakeObjekteFormatterJson() objekteFormatterJson {
	return objekteFormatterJson{}
}

func (f objekteFormatterJson) Format(w io.Writer, o *Objekte) (n int64, err error) {
	enc := json.NewEncoder(w)

	if err = enc.Encode(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
