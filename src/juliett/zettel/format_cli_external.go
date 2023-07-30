package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

// [path.akte_ext@sha]
func MakeCliFormatFD(
	s standort.Standort,
	cw format.FuncColorWriter,
) schnittstellen.FuncWriterFormat[kennung.FD] {
	return func(w io.Writer, fd kennung.FD) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(fd.Path), format.ColorTypePointer),
			// format.MakeFormatString("@"),
			// format.MakeWriter(sf, fd.Sha),
			format.MakeFormatString("]"),
		)
	}
}
