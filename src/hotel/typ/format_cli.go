package typ

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

// !typ
func MakeCliFormat(
	cw format.FuncColorWriter,
) schnittstellen.FuncWriterFormat[kennung.Typ] {
	return func(w io.Writer, t kennung.Typ) (n int64, err error) {
		v := t.String()

		return format.Write(
			w,
			format.MakeFormatString("!"),
			cw(format.MakeFormatString("%s", v), format.ColorTypeType),
		)
	}
}

// [typ.typ@sha !typ]
func MakeCliFormatCheckedOut(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	tf schnittstellen.FuncWriterFormat[kennung.Typ],
) schnittstellen.FuncWriterFormat[CheckedOut] {
	return func(w io.Writer, t CheckedOut) (n int64, err error) {
		diff := format.StringChanged

		if t.Internal.Sku.ObjekteSha.Equals(t.External.Sku.ObjekteSha) {
			diff = format.StringSame
		}

		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(diff),
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(t.External.GetObjekteFD().Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, t.External.GetObjekteSha().GetSha()),
			format.MakeFormatString(" "),
			format.MakeWriter(tf, t.External.Sku.Kennung),
			format.MakeFormatString("]"),
		)
	}
}

// [typ.typ@sha !typ]
func MakeCliFormatTransacted(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	tf schnittstellen.FuncWriterFormat[kennung.Typ],
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
