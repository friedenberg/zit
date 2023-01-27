package zettel_checked_out

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

// (same|changed) [path@sha !typ "bez"]
// (same|changed) [path.akte_ext@sha]
func MakeCliFormat(
	s standort.Standort,
	cw format.FuncColorWriter,
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	zf schnittstellen.FuncWriterFormat[zettel.Objekte],
) schnittstellen.FuncWriterFormat[Zettel] {
	wzef := makeWriterFuncZettel(
		s, cw, hf, sf, zf,
	)

	waef := makeWriterFuncAkte(
		s, cw, hf, sf, zf,
	)

	return func(w io.Writer, z Zettel) (n int64, err error) {
		switch {
		case z.External.AkteFD.Path != "" && z.External.ZettelFD.Path != "":
			return format.Write(
				w,
				format.MakeWriter(wzef, z),
				format.MakeFormatString("\n"),
				format.MakeWriter(waef, z),
			)

		case z.External.AkteFD.Path != "":
			return format.Write(
				w,
				format.MakeWriter(waef, z),
			)

		case z.External.ZettelFD.Path != "":
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
	zf schnittstellen.FuncWriterFormat[zettel.Objekte],
) schnittstellen.FuncWriterFormat[Zettel] {
	return func(w io.Writer, z Zettel) (n int64, err error) {
		diff := format.StringChanged

		if z.Internal.Sku.ObjekteSha.Equals(z.External.Sku.ObjekteSha) {
			diff = format.StringSame
		}

		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(diff),
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(z.External.ZettelFD.Path), format.ColorTypePointer),
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
	zf schnittstellen.FuncWriterFormat[zettel.Objekte],
) schnittstellen.FuncWriterFormat[Zettel] {
	return func(w io.Writer, z Zettel) (n int64, err error) {
		diff := format.StringChanged

		if z.Internal.Objekte.Akte.Equals(z.External.Objekte.Akte) {
			diff = format.StringSame
		}

		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(diff),
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(z.External.AkteFD.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, z.External.Objekte.Akte.GetSha()),
			format.MakeFormatString("]"),
		)
	}
}
