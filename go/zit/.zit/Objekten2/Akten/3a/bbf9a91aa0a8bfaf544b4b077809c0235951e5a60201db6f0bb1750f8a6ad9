package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
)

func (u *Local) FormatOutputOptions() (o string_format_writer.OutputOptions) {
	o.ColorOptionsOut = u.FormatColorOptionsOut()
	o.ColorOptionsErr = u.FormatColorOptionsErr()
	return
}

func (u *Local) FormatColorOptionsOut() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.out.IsTty || !u.config.PrintOptions.PrintColors
	return
}

func (u *Local) FormatColorOptionsErr() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.err.IsTty || !u.config.PrintOptions.PrintColors
	return
}

func (u *Local) StringFormatWriterFields(
	truncate string_format_writer.CliFormatTruncation,
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[string_format_writer.Box] {
	return string_format_writer.MakeCliFormatFields(truncate, co)
}
