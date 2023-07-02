package etikett

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

func MakeCliFormat(
	cw format.FuncColorWriter,
) schnittstellen.FuncWriterFormat[kennung.Etikett] {
	return func(w io.Writer, t kennung.Etikett) (n int64, err error) {
		v := t.String()

		return format.Write(
			w,
			format.MakeFormatString("-"),
			cw(format.MakeFormatString("%s", v), format.ColorTypeType),
		)
	}
}

// [etikett.etikett@sha -etikett]
func MakeCliFormatCheckedOut(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	tf schnittstellen.FuncWriterFormat[kennung.Etikett],
) schnittstellen.FuncWriterFormat[CheckedOut] {
	return func(w io.Writer, t CheckedOut) (n int64, err error) {
		// diff := format.StringChanged

		// if t.Internal.Sku.ObjekteSha.Equals(t.External.Sku.ObjekteSha) {
		// 	diff = format.StringSame
		// }

		return format.Write(
			w,
			// format.MakeFormatStringRightAlignedParen(diff),
			format.MakeFormatString("["),
			cw(
				s.MakeWriterRelativePath(t.External.GetObjekteFD().Path),
				format.ColorTypePointer,
			),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, t.External.GetObjekteSha().GetSha()),
			format.MakeFormatString(" "),
			format.MakeWriter(tf, t.External.Sku.GetKennung()),
			format.MakeFormatString("]"),
		)
	}
}

// [etikett.etikett@sha ]
func MakeCliFormatExternal(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	tf schnittstellen.FuncWriterFormat[kennung.Etikett],
) schnittstellen.FuncWriterFormat[External] {
	return func(w io.Writer, t External) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAligned(""),
			format.MakeFormatString("["),
			cw(
				s.MakeWriterRelativePath(t.GetObjekteFD().Path),
				format.ColorTypePointer,
			),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, t.GetObjekteSha()),
			format.MakeFormatString(" "),
			format.MakeWriter(tf, t.Sku.GetKennung()),
			format.MakeFormatString("]"),
		)
	}
}

// [etikett.etikett@sha ]
func MakeCliFormatTransacted(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	tf schnittstellen.FuncWriterFormat[kennung.Etikett],
) schnittstellen.FuncWriterFormat[Transacted] {
	return func(w io.Writer, t Transacted) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			cw(format.MakeWriter(tf, *t.Kennung()), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, t.GetObjekteSha()),
			format.MakeFormatString("]"),
		)
	}
}
