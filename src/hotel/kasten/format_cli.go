package kasten

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

func MakeCliFormat(
	cw format.FuncColorWriter,
) schnittstellen.FuncWriterFormat[kennung.Kasten] {
	return func(w io.Writer, t kennung.Kasten) (n int64, err error) {
		v := t.String()

		return format.Write(
			w,
			format.MakeFormatString("//"),
			cw(format.MakeFormatString("%s", v), format.ColorTypeType),
		)
	}
}

// [id.kasten@sha]
func MakeCliFormatCheckedOut(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	tf schnittstellen.FuncWriterFormat[kennung.Kasten],
) schnittstellen.FuncWriterFormat[CheckedOut] {
	return func(w io.Writer, t CheckedOut) (n int64, err error) {
		return format.Write(
			w,
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

// [kasten.kasten@sha ]
func MakeCliFormatExternal(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	tf schnittstellen.FuncWriterFormat[kennung.Kasten],
) schnittstellen.FuncWriterFormat[External] {
	return func(w io.Writer, t External) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAligned(""),
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(t.GetObjekteFD().Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, t.GetObjekteSha()),
			format.MakeFormatString(" "),
			format.MakeWriter(tf, t.Sku.Kennung),
			format.MakeFormatString("]"),
		)
	}
}

// [kasten.kasten@sha ]
func MakeCliFormatTransacted(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	tf schnittstellen.FuncWriterFormat[kennung.Kasten],
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
