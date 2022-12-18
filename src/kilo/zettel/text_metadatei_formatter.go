package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/charlie/gattung"
)

//TODO-P1 implement for zettel_external
type TextMetadateiFormatter struct {
	AkteFactory gattung.AkteIOFactory
}

func (f *TextMetadateiFormatter) Format(w1 io.Writer, z *Objekte) (n int64, err error) {
	return
}
