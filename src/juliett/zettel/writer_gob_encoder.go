package zettel

import (
	"encoding/gob"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

// TODO-P1 remove and make transform func in collections
type writerGobEncoder struct {
	enc *gob.Encoder
}

func MakeWriterGobEncoder(w io.Writer) writerGobEncoder {
	return writerGobEncoder{
		enc: gob.NewEncoder(w),
	}
}

func (w writerGobEncoder) WriteZettelVerzeichnisse(z *Transacted) (err error) {
	if err = w.enc.Encode(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
