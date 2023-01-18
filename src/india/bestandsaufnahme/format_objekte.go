package bestandsaufnahme

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/format"
)

type formatObjekte struct {
}

func MakeFormatObjekte() *formatObjekte {
	return &formatObjekte{}
}

func (f *formatObjekte) Parse(r io.Reader, o *Objekte) (n int64, err error) {
	if n, err = format.ReadLines(
		r,
		format.MakeLineReaderKeyValue("Tai", o.Tai.Set),
		format.MakeLineReaderKeyValue("Akte", o.AkteSha.Set),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *formatObjekte) Format(w io.Writer, o *Objekte) (n int64, err error) {
	if n, err = format.WriteLines(
		w,
		format.MakeFormatString("Tai %s", o.Tai),
		format.MakeFormatString("Akte %s", o.AkteSha),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
