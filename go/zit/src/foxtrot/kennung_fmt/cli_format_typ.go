package kennung_fmt

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/src/echo/kennung"
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
	w schnittstellen.WriterAndStringWriter,
	k *kennung.Typ,
) (n int64, err error) {
	v := k.String()

	return f.stringFormatWriter.WriteStringFormat(w, v)
}
