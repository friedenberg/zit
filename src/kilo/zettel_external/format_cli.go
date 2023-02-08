package zettel_external

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

// [path@sha !typ "bez"]
// [path.akte_ext@sha]
func MakeCliFormat(
	s standort.Standort,
	cw format.FuncColorWriter,
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	zf schnittstellen.FuncWriterFormat[zettel.Objekte],
) schnittstellen.FuncWriterFormat[Zettel] {
	return func(w io.Writer, z Zettel) (n int64, err error) {
		switch {
		case z.AkteFD.Path != "" && z.ZettelFD.Path != "":
			return format.Write(
				w,
				format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.ZettelFD.Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.GetObjekteSha().GetSha()),
				format.MakeFormatString(" "),
				format.MakeWriter(zf, z.Objekte),
				format.MakeFormatString("]\n"),
				format.MakeFormatStringRightAlignedParen(""),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.AkteFD.Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.Objekte.Akte.GetSha()),
				format.MakeFormatString("]"),
			)

		case z.AkteFD.Path != "":
			return format.Write(
				w,
				format.MakeFormatStringRightAlignedParen(format.StringCheckedOut),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.AkteFD.Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.Objekte.Akte.GetSha()),
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
				format.MakeWriter(sf, z.GetObjekteSha().GetSha()),
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
