package typ

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/kennung"
)

// !typ
func MakeKennungCliFormat(
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

// !typ
func MakeCliFormat(
	cw format.FuncColorWriter,
) format.FormatWriterFunc[Named] {
	return func(w io.Writer, t *Named) (n int64, err error) {
		v := t.Kennung.String()

		return format.Write(
			w,
			format.MakeFormatString("!"),
			cw(format.MakeFormatString("%s", v), format.ColorTypeType),
		)
	}
}
