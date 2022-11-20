package typ

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/konfig"
)

type EncoderActionNames struct {
	out    io.Writer
	konfig konfig.Compiled
}

func MakeEncoderActionNames(w io.Writer, k konfig.Konfig) EncoderActionNames {
	return EncoderActionNames{
		out:    w,
		konfig: k.Compiled,
	}
}

func (f EncoderActionNames) Encode(t *Typ) (n int64, err error) {
	ct := f.konfig.GetTyp(t.String())

	if ct == nil {
		return
	}

	for v, _ := range ct.Actions {
		var n1 int

		if n1, err = io.WriteString(f.out, fmt.Sprintf("%s\n", v)); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)
	}

	return
}
