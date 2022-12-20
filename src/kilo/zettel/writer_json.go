package zettel

import (
	"encoding/json"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type WriterJson struct {
	enc *json.Encoder
}

func MakeWriterJson(w io.Writer) (w1 WriterJson) {
	return WriterJson{
		enc: json.NewEncoder(w),
	}
}

func (w WriterJson) WriteZettelVerzeichnisse(z *Verzeichnisse) (err error) {
	errors.Log().Printf("writing zettel: %v", z)
	if err = w.enc.Encode(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
