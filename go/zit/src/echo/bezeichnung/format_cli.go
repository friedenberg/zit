package bezeichnung

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/string_format_writer"
)

type bezeichnungCliFormat struct {
	truncate           CliFormatTruncation
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeCliFormat2(
	truncate CliFormatTruncation,
	co string_format_writer.ColorOptions,
	quote bool,
) *bezeichnungCliFormat {
	sfw := string_format_writer.MakeString[string]()

	if quote {
		sfw = string_format_writer.MakeQuotedString[string]()
	}

	return &bezeichnungCliFormat{
		truncate: truncate,
		stringFormatWriter: string_format_writer.MakeColor[string](
			co,
			sfw,
			string_format_writer.ColorTypeIdentifier,
		),
	}
}

func (f *bezeichnungCliFormat) WriteStringFormat(
	w schnittstellen.WriterAndStringWriter,
	k *Bezeichnung,
) (n int64, err error) {
	v := k.value

	if f.truncate == CliFormatTruncation66CharEllipsis && len(v) > 66 {
		v = v[:66] + "…"
	}

	return f.stringFormatWriter.WriteStringFormat(w, v)
}
