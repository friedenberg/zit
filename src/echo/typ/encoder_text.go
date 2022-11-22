package typ

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/metadatei_io"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type EncoderText struct {
	out             io.Writer
	arf             metadatei_io.AkteIOFactory
	IgnoreTypErrors bool
}

func MakeEncoderText(out io.Writer, arf metadatei_io.AkteIOFactory) *EncoderText {
	return &EncoderText{
		out: out,
		arf: arf,
	}
}

func (f EncoderText) Encode(t *Typ) (n int64, err error) {
	var ar sha.ReadCloser

	if ar, err = f.arf.AkteReader(t.Akte.Sha); err != nil {
		err = errors.Wrap(err)
		return
	}

	w := metadatei_io.Writer{
		// Metadatei: ,
		Akte: ar,
	}

	if n, err = w.WriteTo(f.out); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
