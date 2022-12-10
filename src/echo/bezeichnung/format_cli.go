package bezeichnung

import (
	"io"

	"github.com/friedenberg/zit/src/delta/format"
)

func MakeCliFormat(
	cw format.FuncColorWriter,
) format.FormatWriterFunc[Bezeichnung] {
	return func(w io.Writer, b1 *Bezeichnung) (n int64, err error) {
		b := b1.value

		switch {
		case len(b) > 66:
			b = b[:66] + "…"
		}

		return format.Write(
			w,
			format.MakeFormatString("\""),
			cw(format.MakeFormatString("%s", b), format.ColorTypeIdentifier),
			format.MakeFormatString("\""),
		)
	}
}
