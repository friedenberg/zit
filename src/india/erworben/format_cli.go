package erworben

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
)

// (unchanged|updated) [konfig@sha]
func MakeCliFormatTransacted(
	cw format.FuncColorWriter,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	verb string,
) schnittstellen.FuncWriterFormat[Transacted] {
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
