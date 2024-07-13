package descriptions

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
)

type formatCli struct {
	truncate           CliFormatTruncation
	stringFormatWriter interfaces.StringFormatWriter[string]
}

func MakeCliFormat(
	truncate CliFormatTruncation,
	co string_format_writer.ColorOptions,
	quote bool,
) *formatCli {
	sfw := string_format_writer.MakeString[string]()

	if quote {
		sfw = string_format_writer.MakeQuotedString[string]()
	}

	return &formatCli{
		truncate: truncate,
		stringFormatWriter: string_format_writer.MakeColor(
			co,
			sfw,
			string_format_writer.ColorTypeIdentifier,
		),
	}
}

func (f *formatCli) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	k *Description,
) (n int64, err error) {
	v := k.value

	if f.truncate == CliFormatTruncation66CharEllipsis && len(v) > 66 {
		v = v[:66] + "…"
	}

	return f.stringFormatWriter.WriteStringFormat(w, v)
}
