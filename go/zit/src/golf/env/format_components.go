package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
)

func (u *Env) FormatOutputOptions() (o string_format_writer.OutputOptions) {
	o.ColorOptionsOut = u.FormatColorOptionsOut()
	o.ColorOptionsErr = u.FormatColorOptionsErr()
	return
}

func (u *Env) FormatColorOptionsOut() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.GetOut().IsTty() || !u.cliConfig.PrintOptions.PrintColors
	return
}

func (u *Env) FormatColorOptionsErr() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.GetErr().IsTty() || !u.cliConfig.PrintOptions.PrintColors
	return
}

func (u *Env) StringFormatWriterFields(
	truncate string_format_writer.CliFormatTruncation,
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[string_format_writer.Box] {
	return string_format_writer.MakeCliFormatFields(truncate, co)
}
