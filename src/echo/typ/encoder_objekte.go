package typ

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/line_format"
)

type EncoderObjekte struct {
	out             io.Writer
	IgnoreTypErrors bool
}

func MakeEncoderObjekte(out io.Writer) *EncoderObjekte {
	return &EncoderObjekte{
		out: out,
	}
}

func (f EncoderObjekte) Encode(t *Typ) (n int64, err error) {
	w := line_format.NewWriter()

	w.WriteFormat("%s", gattung.Typ)
	w.WriteFormat("%s %s", gattung.Akte, t.Akte.Sha)

	if n, err = w.WriteTo(f.out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
