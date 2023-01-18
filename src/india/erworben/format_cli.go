package erworben

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/echo/format"
)

// (unchanged|updated) [konfig@sha]
func MakeCliFormatTransacted(
	cw format.FuncColorWriter,
	sf format.FormatWriterFunc[sha.Sha],
	verb string,
) format.FormatWriterFunc[Transacted] {
	return func(w io.Writer, kt Transacted) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(verb),
			format.MakeFormatString("["),
			cw(format.MakeFormatString("%s", kt.Kennung()), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, kt.GetObjekteSha()),
			format.MakeFormatString("]"),
		)
	}
}
