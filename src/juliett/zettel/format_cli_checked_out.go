package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

// (same|changed) [path@sha !typ "bez"]
// (same|changed) [path.akte_ext@sha]
func MakeCliFormatCheckedOut(
	s standort.Standort,
	cw format.FuncColorWriter,
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	zf schnittstellen.FuncWriterFormat[Objekte],
) schnittstellen.FuncWriterFormat[CheckedOut] {
	wzef := makeWriterFuncZettel(
		s, cw, hf, sf, zf,
	)

	waef := makeWriterFuncAkte(
		s, cw, hf, sf, zf,
	)

	return func(w io.Writer, z CheckedOut) (n int64, err error) {
		switch {
		case z.External.Sku.FDs.Akte.Path != "" && z.External.Sku.FDs.Objekte.Path != "":
			return format.Write(
				w,
				format.MakeWriter(wzef, z),
				format.MakeFormatString("\n"),
				format.MakeWriter(waef, z),
			)

		case z.External.Sku.FDs.Akte.Path != "":
			return format.Write(
				w,
				format.MakeWriter(waef, z),
			)

		case z.External.Sku.FDs.Objekte.Path != "":
			return format.Write(
				w,
				format.MakeWriter(wzef, z),
			)
		}

		return
	}
}

func makeWriterFuncZettel(
	s standort.Standort,
	cw format.FuncColorWriter,
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	zf schnittstellen.FuncWriterFormat[Objekte],
) schnittstellen.FuncWriterFormat[CheckedOut] {
	return func(w io.Writer, z CheckedOut) (n int64, err error) {
		diff := format.StringChanged

		if z.Internal.Sku.ObjekteSha.Equals(z.External.Sku.ObjekteSha) {
			diff = format.StringSame
		}

		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(diff),
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(z.External.GetObjekteFD().Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, z.External.GetObjekteSha().GetSha()),
			format.MakeFormatString(" "),
			format.MakeWriter(zf, z.External.Objekte),
			format.MakeFormatString("]"),
		)
	}
}

func makeWriterFuncAkte(
	s standort.Standort,
	cw format.FuncColorWriter,
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	zf schnittstellen.FuncWriterFormat[Objekte],
) schnittstellen.FuncWriterFormat[CheckedOut] {
	return func(w io.Writer, z CheckedOut) (n int64, err error) {
		diff := format.StringChanged

		if z.Internal.Objekte.Akte.Equals(z.External.Objekte.Akte) {
			diff = format.StringSame
		}

		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(diff),
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(z.External.GetAkteFD().Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, z.External.Objekte.Akte.GetSha()),
			format.MakeFormatString("]"),
		)
	}
}
