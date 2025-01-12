package read_write_repo_local

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
)

func (u *Repo) FormatOutputOptions() (o string_format_writer.OutputOptions) {
	o.ColorOptionsOut = u.FormatColorOptionsOut()
	o.ColorOptionsErr = u.FormatColorOptionsErr()
	return
}

func (u *Repo) FormatColorOptionsOut() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.GetOut().IsTty() || !u.config.PrintOptions.PrintColors
	return
}

func (u *Repo) FormatColorOptionsErr() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.GetErr().IsTty() || !u.config.PrintOptions.PrintColors
	return
}

func (u *Repo) StringFormatWriterFields(
	truncate string_format_writer.CliFormatTruncation,
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[string_format_writer.Box] {
	return string_format_writer.MakeCliFormatFields(truncate, co)
}
