package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

// [path@sha !typ "bez"]
// [path.akte_ext@sha]
func MakeCliFormat(
	s standort.Standort,
	cw format.FuncColorWriter,
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	mf schnittstellen.FuncWriterFormat[metadatei.GetterPtr],
) schnittstellen.FuncWriterFormat[External] {
	return func(w io.Writer, z External) (n int64, err error) {
		switch {
		case z.GetAkteFD().Path != "" && z.GetObjekteFD().Path != "":
			return format.Write(
				w,
				format.MakeFormatStringRightAligned(format.StringCheckedOut),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.GetObjekteFD().Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.GetObjekteSha().GetSha()),
				format.MakeFormatString(" "),
				format.MakeWriter[metadatei.GetterPtr](mf, &z),
				format.MakeFormatString("]\n"),
				format.MakeFormatStringRightAligned(""),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.GetAkteFD().Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, schnittstellen.Sha(z.GetAkteSha())),
				format.MakeFormatString("]"),
			)

		case z.GetAkteFD().Path != "":
			return format.Write(
				w,
				format.MakeFormatStringRightAligned(format.StringCheckedOut),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.GetAkteFD().Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, schnittstellen.Sha(z.GetAkteSha())),
				format.MakeFormatString(" "),
				format.MakeWriter[metadatei.GetterPtr](mf, &z),
				format.MakeFormatString("]"),
			)

		case z.GetObjekteFD().Path != "":
			return format.Write(
				w,
				format.MakeFormatStringRightAligned(format.StringCheckedOut),
				format.MakeFormatString("["),
				cw(s.MakeWriterRelativePath(z.GetObjekteFD().Path), format.ColorTypePointer),
				format.MakeFormatString("@"),
				format.MakeWriter(sf, z.GetObjekteSha().GetSha()),
				format.MakeFormatString(" "),
				format.MakeWriter[metadatei.GetterPtr](mf, &z),
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
