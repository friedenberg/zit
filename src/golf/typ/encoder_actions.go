package typ

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/konfig"
)

type EncoderActionNames struct {
	out    io.Writer
	konfig konfig.Objekte
}

func MakeEncoderActionNames(w io.Writer, k konfig.Konfig) EncoderActionNames {
	return EncoderActionNames{
		out:    w,
		konfig: k.Transacted.Objekte,
	}
}

func (f EncoderActionNames) Encode(t *kennung.Typ) (n int64, err error) {
	ct := f.konfig.GetTyp(t.String())

	if ct == nil {
		return
	}

	for v, v1 := range ct.Typ.Actions {
		var n1 int

		if n1, err = io.WriteString(f.out, fmt.Sprintf("%s\t%s\n", v, v1.Description)); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)
	}

	return
}
