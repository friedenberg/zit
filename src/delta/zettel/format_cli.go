package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/bezeichnung"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/typ"
)

// !typ "bez"
func MakeCliFormat(
	bf collections.WriterFuncFormat[bezeichnung.Bezeichnung],
	ef collections.WriterFuncFormat[etikett.Set],
	tf collections.WriterFuncFormat[typ.Typ],
) collections.WriterFuncFormat[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		return collections.WriteFormats(
			w,
			collections.MakeWriterLiteral("!"),
			collections.MakeWriterFormatFunc(tf, &z.Typ),
			collections.MakeWriterLiteral(" "),
			collections.MakeWriterFormatFunc(bf, &z.Bezeichnung),
		)
	}
}
