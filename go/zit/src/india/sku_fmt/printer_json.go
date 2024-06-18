package sku_fmt

import (
	"encoding/json"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type WriterJson struct {
	enc *json.Encoder
}

func MakeWriterJson(w io.Writer) (w1 WriterJson) {
	return WriterJson{
		enc: json.NewEncoder(w),
	}
}

func (w WriterJson) WriteZettelVerzeichnisse(z *sku.Transacted) (err error) {
	if err = w.enc.Encode(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
