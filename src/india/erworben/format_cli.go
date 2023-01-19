package erworben

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
)

// (unchanged|updated) [konfig@sha]
func MakeCliFormatTransacted(
	cw format.FuncColorWriter,
	sf format.FormatWriterFunc[schnittstellen.Sha],
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
