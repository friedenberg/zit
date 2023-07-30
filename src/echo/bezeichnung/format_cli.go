package bezeichnung

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/string_writer_format"
)

type bezeichnungCliFormat struct {
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeCliFormat2(
	co string_writer_format.ColorOptions,
) *bezeichnungCliFormat {
	return &bezeichnungCliFormat{
		stringFormatWriter: string_writer_format.MakeColor[string](
			co,
			string_writer_format.MakeString[string](),
			string_writer_format.ColorTypeIdentifier,
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
