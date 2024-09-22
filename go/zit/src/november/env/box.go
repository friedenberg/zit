package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

func (u *Env) StringFormatWriterSkuBox(
	po erworben_cli_print_options.PrintOptions,
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

func (u *Env) SkuFormatBox() sku_fmt.ExternalLike {
	co := u.FormatColorOptionsOut()

	return u.StringFormatWriterSkuBox(
		u.config.PrintOptions.WithPrintShas(false),
		co,
		string_format_writer.CliFormatTruncationNone,
	)
}

func (u *Env) StringFormatWriterSkuTransacted(
	co *string_format_writer.ColorOptions,
	truncate string_format_writer.CliFormatTruncation,
) interfaces.StringFormatWriter[*sku.Transacted] {
	if co == nil {
		co1 := u.FormatColorOptionsOut()
		co = &co1
	}

	return sku_fmt.MakeCliFormat(
		u.config.PrintOptions,
		u.StringFormatWriterObjectId(*co),
		u.StringFormatWriterMetadata(*co, truncate),
	)
}

func (u *Env) StringFormatWriterSkuTransactedShort() interfaces.StringFormatWriter[*sku.Transacted] {
	co := string_format_writer.ColorOptions{
		OffEntirely: true,
	}

	return sku_fmt.MakeCliFormatShort(
		u.StringFormatWriterObjectId(co),
		u.StringFormatWriterMetadata(
			co,
			string_format_writer.CliFormatTruncation66CharEllipsis,
		),
	)
}
