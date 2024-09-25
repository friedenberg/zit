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
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterObjectIdAligned(co),
		u.StringFormatWriterType(co),
		u.StringFormatWriterTags(co),
		u.StringFormatWriterFields(truncation, co),
		u.StringFormatWriterMetadata(
			co,
			truncation,
		),
		u.GetStore().GetAbbrStore().GetAbbr(),
	)
}

func (u *Env) SkuFormatBoxNoColor() sku_fmt.ExternalLike {
	co := u.FormatColorOptionsOut()
  co.OffEntirely = true

	return u.StringFormatWriterSkuBox(
		u.config.PrintOptions.WithPrintShas(false),
		co,
		string_format_writer.CliFormatTruncationNone,
	)
}
