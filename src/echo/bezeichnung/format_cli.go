package bezeichnung

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
)

type bezeichnungCliFormat struct {
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeCliFormat2(
	co string_format_writer.ColorOptions,
) *bezeichnungCliFormat {
	return &bezeichnungCliFormat{
		stringFormatWriter: string_format_writer.MakeColor[string](
			co,
			string_format_writer.MakeString[string](),
			string_format_writer.ColorTypeIdentifier,
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
