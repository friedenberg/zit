package kennung_fmt

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type typCliFormat struct {
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeTypCliFormat(co string_format_writer.ColorOptions) *typCliFormat {
	return &typCliFormat{
		stringFormatWriter: string_format_writer.MakeColor[string](
			co,
			string_format_writer.MakeString[string](),
			string_format_writer.ColorTypeType,
		),
	}
}

func (f *typCliFormat) WriteStringFormat(
	w io.StringWriter,
	k *kennung.Typ,
) (n int64, err error) {
	v := k.String()

	return f.stringFormatWriter.WriteStringFormat(w, v)
}
