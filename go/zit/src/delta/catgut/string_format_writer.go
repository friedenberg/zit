package catgut

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type stringFormatWriter struct{}

var StringFormatWriterString stringFormatWriter

func (stringFormatWriter) EncodeStringTo(
	e *String,
	sw interfaces.WriterAndStringWriter,
) (n int64, err error) {
	n, err = e.WriteTo(sw)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
