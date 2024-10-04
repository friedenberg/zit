package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/id_fmts"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata_fmt"
)

func (u *Env) FormatOutputOptions() (o string_format_writer.OutputOptions) {
	o.ColorOptionsOut = u.FormatColorOptionsOut()
	o.ColorOptionsErr = u.FormatColorOptionsErr()
	return
}

func (u *Env) FormatColorOptionsOut() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.outIsTty || !u.config.PrintOptions.PrintColors
	return
}

func (u *Env) FormatColorOptionsErr() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.errIsTty || !u.config.PrintOptions.PrintColors
	return
}

func (u *Env) StringFormatWriterShaLike(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[interfaces.Sha] {
	return id_fmts.MakeShaCliFormat(
		u.config.PrintOptions,
		co,
		u.store.GetAbbrStore().GetAbbr(),
	)
}

func (u *Env) StringFormatWriterObjectIdAligned(
	co string_format_writer.ColorOptions,
) id_fmts.Aligned {
	return id_fmts.MakeAligned(
		u.config.PrintOptions,
		u.GetStore().GetAbbrStore().GetAbbr(),
	)
}

func (u *Env) StringFormatWriterObjectId(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*ids.ObjectId] {
	return id_fmts.MakeObjectIdCliFormat(
		u.config.PrintOptions,
		co,
		u.GetStore().GetAbbrStore().GetAbbr(),
	)
}

func (u *Env) StringFormatWriterType(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*ids.Type] {
	return id_fmts.MakeTypCliFormat(co)
}

func (u *Env) StringFormatWriterDescription(
	truncate string_format_writer.CliFormatTruncation,
	co string_format_writer.ColorOptions,
	quote bool,
) interfaces.StringFormatWriter[*descriptions.Description] {
	return descriptions.MakeCliFormat(truncate, co, quote)
}

func (u *Env) StringFormatWriterFields(
	truncate string_format_writer.CliFormatTruncation,
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[string_format_writer.Box] {
	return string_format_writer.MakeCliFormatFields(truncate, co)
}

func (u *Env) StringFormatWriterTags(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*ids.Tag] {
	return id_fmts.MakeEtikettenCliFormat()
}

func (u *Env) StringFormatWriterMetadata(
	co string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
) *object_metadata_fmt.Box {
	return object_metadata_fmt.MakeBoxMetadataFormat(
		u.config.PrintOptions,
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterType(co),
		u.StringFormatWriterFields(truncation, co),
		u.StringFormatWriterTags(co),
	)
}
