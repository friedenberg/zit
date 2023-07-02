package zettel

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

// (same|changed) [path@sha !typ "bez"]
// (same|changed) [path.akte_ext@sha]
func MakeCliFormatCheckedOut(
	s standort.Standort,
	cw format.FuncColorWriter,
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.ShaLike],
	mf schnittstellen.FuncWriterFormat[metadatei.GetterPtr],
) schnittstellen.FuncWriterFormat[CheckedOut] {
	wzef := makeWriterFuncZettel(
		s, cw, hf, sf, mf,
	)

	waef := makeWriterFuncAkte(
		s, cw, hf, sf, mf,
	)

	return func(w io.Writer, z CheckedOut) (n int64, err error) {
		if z.External.Sku.FDs.Akte.Path == "" {
			return format.Write(
				w,
				format.MakeWriter(wzef, z),
			)
		} else {
			return format.Write(
				w,
				format.MakeWriter(wzef, z),
				format.MakeFormatString("\n"),
				format.MakeWriter(waef, z),
			)
		}
	}
}

func makeWriterFuncZettel(
	s standort.Standort,
	cw format.FuncColorWriter,
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.ShaLike],
	mf schnittstellen.FuncWriterFormat[metadatei.GetterPtr],
) schnittstellen.FuncWriterFormat[CheckedOut] {
	return func(w io.Writer, z CheckedOut) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			cw(
				s.MakeWriterRelativePathOr(
					z.External.GetObjekteFD().Path,
					format.MakeWriter(hf, z.External.Sku.GetKennung()),
				),
				format.ColorTypePointer,
			),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, z.External.GetObjekteSha().GetShaLike()),
			format.MakeFormatString(" "),
			format.MakeWriter[metadatei.GetterPtr](mf, &z.External),
			format.MakeFormatString("]"),
		)
	}
}

func makeWriterFuncAkte(
	s standort.Standort,
	cw format.FuncColorWriter,
	hf schnittstellen.FuncWriterFormat[kennung.Hinweis],
	sf schnittstellen.FuncWriterFormat[schnittstellen.ShaLike],
	mf schnittstellen.FuncWriterFormat[metadatei.GetterPtr],
) schnittstellen.FuncWriterFormat[CheckedOut] {
	return func(w io.Writer, z CheckedOut) (n int64, err error) {
		todo.Change("refactor to support proper spacing")
		return format.Write(
			w,
			format.MakeFormatStringRightAligned("%s", format.StringDRArrow),
			format.MakeFormatString("["),
			cw(
				s.MakeWriterRelativePath(z.External.GetAkteFD().Path),
				format.ColorTypePointer,
			),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, z.External.GetAkteSha().GetShaLike()),
			format.MakeFormatString("]"),
		)
	}
}
