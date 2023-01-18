package typ

import (
	"io"

	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/standort"
	"github.com/friedenberg/zit/src/schnittstellen"
)

// !typ
func MakeCliFormat(
	cw format.FuncColorWriter,
) format.FormatWriterFunc[kennung.Typ] {
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
	sf format.FormatWriterFunc[schnittstellen.Sha],
	tf format.FormatWriterFunc[kennung.Typ],
) format.FormatWriterFunc[External] {
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
	sf format.FormatWriterFunc[schnittstellen.Sha],
	tf format.FormatWriterFunc[kennung.Typ],
	verb string,
) format.FormatWriterFunc[Transacted] {
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
