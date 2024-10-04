package env

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

func (u *Env) StringFormatWriterSkuBox(
	po print_options.General,
	co string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
) *sku_fmt.Box {
	return sku_fmt.MakeBox(
		co,
		po,
		u.StringFormatWriterFields(truncation, co),
		u.GetStore().GetAbbrStore().GetAbbr(),
		u.GetStore().GetCwdFiles(),
		u.fs_home,
	)
}

func (u *Env) SkuFormatBoxNoColor() *sku_fmt.Box {
	co := u.FormatColorOptionsOut()
	co.OffEntirely = true
	options := u.config.PrintOptions.WithPrintShas(false)
	options.PrintTime = false
	options.PrintShas = false

	return u.StringFormatWriterSkuBox(
		options,
		co,
		string_format_writer.CliFormatTruncationNone,
	)
}
