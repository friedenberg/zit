package kennung_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type typCliFormat struct {
	stringFormatWriter interfaces.StringFormatWriter[string]
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
	w interfaces.WriterAndStringWriter,
	k *ids.Type,
) (n int64, err error) {
	v := k.String()

	return f.stringFormatWriter.WriteStringFormat(w, v)
}
