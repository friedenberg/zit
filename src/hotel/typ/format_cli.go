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
func MakeCliFormatExternal(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	tf schnittstellen.FuncWriterFormat[kennung.Typ],
) schnittstellen.FuncWriterFormat[External] {
	return func(w io.Writer, t External) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(""),
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(t.FD.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, t.GetObjekteSha().GetSha()),
			format.MakeFormatString(" "),
			format.MakeWriter(tf, t.Sku.Kennung),
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
	verb string,
) schnittstellen.FuncWriterFormat[Transacted] {
	return func(w io.Writer, t Transacted) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(verb),
			format.MakeFormatString("["),
			cw(format.MakeWriter(tf, *t.Kennung()), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, t.GetObjekteSha()),
			format.MakeFormatString("]"),
		)
	}
}
