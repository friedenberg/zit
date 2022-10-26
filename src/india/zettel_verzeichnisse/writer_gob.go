package zettel_verzeichnisse

import (
	"encoding/gob"
	"io"
)

type writerGobEncoder struct {
	enc *gob.Encoder
}

func MakeWriterGobEncoder(w io.Writer) Writer {
	return writerGobEncoder{
		enc: gob.NewEncoder(w),
	}
}

func (w writerGobEncoder) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	return w.enc.Encode(z)
}

type writerGobDecoder struct {
	dec *gob.Decoder
}

func MakeWriterGobDecoder(r io.Reader) Writer {
	return writerGobDecoder{
		dec: gob.NewDecoder(r),
	}
}

func (w writerGobDecoder) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	return w.dec.Decode(z)
}
