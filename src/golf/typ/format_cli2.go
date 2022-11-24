package typ

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
)

// [typ.typ@sha !typ]
func MakeCliFormatTyp(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf format.FormatWriterFunc[sha.Sha],
	tf format.FormatWriterFunc[Named],
) format.FormatWriterFunc[Typ] {
	return func(w io.Writer, t *Typ) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(""),
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(t.FD.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &t.Stored.Sha),
			format.MakeFormatString(" "),
			format.MakeWriter(tf, &t.Named),
			format.MakeFormatString("]"),
		)
	}
}
