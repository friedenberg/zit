package env

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
)

func (u *Env) StringFormatWriterSkuBox(
	po print_options.General,
	co string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
) *box_format.Box {
	return box_format.MakeBox(
		co,
		po,
		u.StringFormatWriterFields(truncation, co),
		u.GetStore().GetAbbrStore().GetAbbr(),
		u.GetStore().GetCwdFiles(),
		u.fsHome,
	)
}

func (u *Env) SkuFormatBoxNoColor() *box_format.Box {
	co := u.FormatColorOptionsOut()
	co.OffEntirely = true
	options := u.config.PrintOptions.WithPrintShas(false)
	options.PrintTime = false
	options.PrintShas = false
	options.DescriptionInBox = false

	return u.StringFormatWriterSkuBox(
		options,
		co,
		string_format_writer.CliFormatTruncationNone,
	)
}

func (u *Env) MakeBoxArchive(includeTai bool) *box_format.Box {
	po := u.GetConfig().PrintOptions.
		WithPrintShas(true).
		WithPrintTai(includeTai).
		WithExcludeFields(true).
		WithDescriptionInBox(true)

	co := u.FormatColorOptionsOut()
	co.OffEntirely = true

	return box_format.MakeBox(
		co,
		po,
		u.StringFormatWriterFields(
			string_format_writer.CliFormatTruncationNone,
			co,
		),
		ids.Abbr{},
		nil,
		nil,
	)
}
