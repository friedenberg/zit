package zettel_verzeichnisse

import (
	"encoding/gob"
	"io"
)

type writerGobEncoder struct {
	enc *gob.Encoder
}

func MakeWriterGobEncoder(w io.Writer) writerGobEncoder {
	return writerGobEncoder{
		enc: gob.NewEncoder(w),
	}
}

func (w writerGobEncoder) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	return w.enc.Encode(z)
}
