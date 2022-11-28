package typ

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/standort"
)

// !typ
func MakeCliFormat(
	cw format.FuncColorWriter,
) format.FormatWriterFunc[kennung.Typ] {
	return func(w io.Writer, t *kennung.Typ) (n int64, err error) {
		v := t.String()

		return format.Write(
			w,
			format.MakeFormatString("!"),
			cw(format.MakeFormatString("%s", v), format.ColorTypeType),
		)
	}
}

// [typ.typ@sha !typ]
func MakeCliFormatExternal(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf format.FormatWriterFunc[sha.Sha],
	tf format.FormatWriterFunc[kennung.Typ],
) format.FormatWriterFunc[External] {
	return func(w io.Writer, t *External) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(""),
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(t.FD.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &t.Sha),
			format.MakeFormatString(" "),
			format.MakeWriter(tf, &t.Kennung),
			format.MakeFormatString("]"),
		)
	}
}
