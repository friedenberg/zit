package zettel_external

import (
	"io"

	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/golf/fd"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

// [path@sha !typ "bez"]
// [path.akte_ext@sha]
func MakeCliFormat(
	s standort.Standort,
	cw format.FuncColorWriter,
	hf format.FormatWriterFunc[hinweis.Hinweis],
	sf format.FormatWriterFunc[sha.Sha],
	zf format.FormatWriterFunc[zettel.Objekte],
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z Zettel) (n int64, err error) {
		switch {
		case z.AkteFD.Path != "" && z.ZettelFD.Path != "":
			return format.Write(
				w,
				format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.ZettelFD.Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.GetObjekteSha()),
				format.MakeFormatString(" "),
				format.MakeWriter(zf, z.Objekte),
				format.MakeFormatString("]\n"),
				format.MakeFormatStringRightAlignedParen(""),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.AkteFD.Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.Objekte.Akte),
				format.MakeFormatString("]"),
			)

		case z.AkteFD.Path != "":
			return format.Write(
				w,
				format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.AkteFD.Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.Objekte.Akte),
				format.MakeFormatString(" "),
				format.MakeWriter(zf, z.Objekte),
				format.MakeFormatString("]"),
			)

		case z.ZettelFD.Path != "":
			return format.Write(
				w,
				format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.ZettelFD.Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.GetObjekteSha()),
				format.MakeFormatString(" "),
				format.MakeWriter(zf, z.Objekte),
				format.MakeFormatString("]"),
			)
		}

		return
	}
}

// [path.akte_ext@sha]
func MakeCliFormatFD(
	s standort.Standort,
	cw format.FuncColorWriter,
) format.FormatWriterFunc[fd.FD] {
	return func(w io.Writer, fd fd.FD) (n int64, err error) {
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
