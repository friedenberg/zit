package metadatei_io

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type WriterToFromReader struct {
	io.Reader
}

func (wtfr WriterToFromReader) WriteTo(w io.Writer) (n int64, err error) {
	if n, err = io.Copy(w, wtfr.Reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
