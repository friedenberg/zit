package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/bezeichnung"
	"github.com/friedenberg/zit/src/delta/etikett"
	"github.com/friedenberg/zit/src/echo/typ"
)

// !typ "bez"
func MakeCliFormat(
	bf format.FormatWriterFunc[bezeichnung.Bezeichnung],
	ef format.FormatWriterFunc[etikett.Set],
	tf format.FormatWriterFunc[typ.Typ],
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		var lastWriter format.WriterFunc

		if z.Bezeichnung.IsEmpty() {
			lastWriter = format.MakeWriter(ef, &z.Etiketten)
		} else {
			lastWriter = format.MakeWriter(bf, &z.Bezeichnung)
		}

		return format.Write(
			w,
			format.MakeWriter(tf, &z.Typ),
			format.MakeFormatString(" "),
			lastWriter,
		)
	}
}