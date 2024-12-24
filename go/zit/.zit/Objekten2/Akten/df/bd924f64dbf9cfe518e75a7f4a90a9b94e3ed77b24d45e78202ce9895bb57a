package descriptions

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
)

type formatCli[T interfaces.Stringer] struct {
	*formatCliStringer
}

func MakeCliFormat(
	truncate string_format_writer.CliFormatTruncation,
	co string_format_writer.ColorOptions,
	quote bool,
) *formatCli[*Description] {
	return MakeCliFormatGeneric[*Description](
		truncate,
		co,
		quote,
	)
}

func MakeCliFormatGeneric[T interfaces.Stringer](
	truncate string_format_writer.CliFormatTruncation,
	co string_format_writer.ColorOptions,
	quote bool,
) *formatCli[T] {
	return &formatCli[T]{
		formatCliStringer: MakeCliFormatStringer(
			truncate,
			co,
			quote,
		),
	}
}

func (f *formatCli[T]) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	k T,
) (n int64, err error) {
	return f.formatCliStringer.WriteStringFormat(w, k)
}
