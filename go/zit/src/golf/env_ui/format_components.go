package env_ui

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
)

func (u *env) FormatOutputOptions() (o string_format_writer.OutputOptions) {
	o.ColorOptionsOut = u.FormatColorOptionsOut()
	o.ColorOptionsErr = u.FormatColorOptionsErr()
	return
}

func (u *env) FormatColorOptionsOut() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.GetOut().IsTty() || !u.cliConfig.PrintOptions.PrintColors
	return
}

func (u *env) FormatColorOptionsErr() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.GetErr().IsTty() || !u.cliConfig.PrintOptions.PrintColors
	return
}

func (u *env) StringFormatWriterFields(
	truncate string_format_writer.CliFormatTruncation,
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[string_format_writer.Box] {
	return string_format_writer.MakeCliFormatFields(truncate, co)
}
