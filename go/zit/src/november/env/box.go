package env

import (
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
)

func (u *Env) StringFormatWriterSkuBoxTransacted(
	po options_print.V0,
	co string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
) *box_format.BoxTransacted {
	return box_format.MakeBoxTransacted(
		co,
		po,
		u.StringFormatWriterFields(truncation, co),
		u.GetStore().GetAbbrStore().GetAbbr(),
		u.GetStore().GetStoreFS(),
		u.dirLayout,
	)
}

func (u *Env) StringFormatWriterSkuBoxCheckedOut(
	po options_print.V0,
	co string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
) *box_format.BoxCheckedOut {
	return box_format.MakeBoxCheckedOut(
		co,
		po,
		u.StringFormatWriterFields(truncation, co),
		u.GetStore().GetAbbrStore().GetAbbr(),
		u.GetStore().GetStoreFS(),
		u.dirLayout,
	)
}

func (u *Env) SkuFormatBoxNoColor() *box_format.BoxTransacted {
	co := u.FormatColorOptionsOut()
	co.OffEntirely = true
	options := u.config.PrintOptions.WithPrintShas(false)
	options.PrintTime = false
	options.PrintShas = false
	options.DescriptionInBox = false

	return u.StringFormatWriterSkuBoxTransacted(
		options,
		co,
		string_format_writer.CliFormatTruncationNone,
	)
}

func (u *Env) MakeBoxArchive(includeTai bool) *box_format.BoxTransacted {
	po := u.GetConfig().PrintOptions.
		WithPrintShas(true).
		WithPrintTai(includeTai).
		WithExcludeFields(true).
		WithDescriptionInBox(true)

	co := u.FormatColorOptionsOut()
	co.OffEntirely = true

	return box_format.MakeBoxTransacted(
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
