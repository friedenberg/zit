package bezeichnung

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
)

func MakeCliFormat(
	cw format.FuncColorWriter,
) schnittstellen.FuncWriterFormat[Bezeichnung] {
	return func(w io.Writer, b1 Bezeichnung) (n int64, err error) {
		b := b1.value

		switch {
		case len(b) > 66:
			b = b[:66] + "…"
		}

		return format.Write(
			w,
			format.MakeFormatString("\""),
			cw(format.MakeFormatString("%s", b), format.ColorTypeIdentifier),
			format.MakeFormatString("\""),
		)
	}
}

type bezeichnungCliFormat struct {
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeCliFormat2(co format.ColorOptions) *bezeichnungCliFormat {
	return &bezeichnungCliFormat{
		stringFormatWriter: format.MakeColorStringFormatWriter[string](
			co,
			format.MakeStringStringFormatWriter[string](),
			format.ColorTypeIdentifier,
		),
	}
}

func (f *bezeichnungCliFormat) WriteStringFormat(
	w io.StringWriter,
	k *Bezeichnung,
) (n int64, err error) {
	v := k.value

	switch {
	case len(v) > 66:
		v = v[:66] + "…"
	}

	return f.stringFormatWriter.WriteStringFormat(w, v)
}
