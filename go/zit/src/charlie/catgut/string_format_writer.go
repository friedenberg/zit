package catgut

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type stringFormatWriter struct{}

var StringFormatWriter stringFormatWriter

func (stringFormatWriter) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	e *String,
) (n int64, err error) {
	n, err = e.WriteTo(sw)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
