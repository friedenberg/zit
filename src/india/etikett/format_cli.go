package etikett

import (
	"io"

	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
)

func MakeCliFormat(
	cw format.FuncColorWriter,
) format.FormatWriterFunc[kennung.Etikett] {
	return func(w io.Writer, t *kennung.Etikett) (n int64, err error) {
		v := t.String()

		return format.Write(
			w,
			format.MakeFormatString("!"),
			cw(format.MakeFormatString("%s", v), format.ColorTypeType),
		)
	}
}

// [etikett.etikett@sha ]
func MakeCliFormatExternal(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf format.FormatWriterFunc[sha.Sha],
	tf format.FormatWriterFunc[kennung.Etikett],
) format.FormatWriterFunc[External] {
	return func(w io.Writer, t *External) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(""),
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(t.FD.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &t.Sku.ObjekteSha),
			format.MakeFormatString(" "),
			format.MakeWriter(tf, &t.Sku.Kennung),
			format.MakeFormatString("]"),
		)
	}
}

// [etikett.etikett@sha ]
func MakeCliFormatTransacted(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf format.FormatWriterFunc[sha.Sha],
	tf format.FormatWriterFunc[kennung.Etikett],
	verb string,
) format.FormatWriterFunc[Transacted] {
	return func(w io.Writer, t *Transacted) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(verb),
			format.MakeFormatString("["),
			cw(format.MakeWriter(tf, t.Kennung()), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &t.Sku.ObjekteSha),
			format.MakeFormatString("]"),
		)
	}
}
