package typ

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
)

// !typ
func MakeCliFormat(
	cw format.FuncColorWriter,
) format.FormatWriterFunc[Kennung] {
	return func(w io.Writer, t *Kennung) (n int64, err error) {
		v := t.String()

		return format.Write(
			w,
			format.MakeFormatString("!"),
			cw(format.MakeFormatString("%s", v), format.ColorTypeType),
		)
	}
}
