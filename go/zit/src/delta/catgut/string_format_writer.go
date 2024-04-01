package catgut

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
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
