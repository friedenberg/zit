package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/id_fmts"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
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
		func(s *sha.Sha) (string, error) {
			return u.store.GetAbbrStore().Shas().Abbreviate(s)
		},
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

func (u *Env) StringFormatWriterTyp(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*ids.Type] {
	return id_fmts.MakeTypCliFormat(co)
}

func (u *Env) StringFormatWriterDescription(
	truncate descriptions.CliFormatTruncation,
	co string_format_writer.ColorOptions,
	quote bool,
) interfaces.StringFormatWriter[*descriptions.Description] {
	return descriptions.MakeCliFormat(truncate, co, quote)
}

func (u *Env) StringFormatWriterEtiketten(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*ids.Tag] {
	return id_fmts.MakeEtikettenCliFormat()
}

func (u *Env) StringFormatWriterMetadatei(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*object_metadata.Metadata] {
	return sku_fmt.MakeCliMetadateiFormat(
		u.config.PrintOptions,
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterDescription(
			descriptions.CliFormatTruncation66CharEllipsis,
			co,
			true,
		),
		u.StringFormatWriterEtiketten(co),
	)
}

func (u *Env) SkuFmtOrganize(repoId ids.RepoId) sku_fmt.ExternalLikeFormatter {
	co := u.FormatColorOptionsOut()
	co.OffEntirely = true

	f := sku_fmt.MakeFormatOrganize(
		u.config.PrintOptions,
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterObjectIdAligned(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterDescription(descriptions.CliFormatTruncationNone, co, false),
		u.StringFormatWriterEtiketten(co),
	)

	kid := repoId.GetRepoIdString()
	es, ok := u.externalStores[kid]

	if !ok {
		return f
	}

	return es.GetExternalStoreOrganizeFormat(f)
}

func (u *Env) StringFormatWriterSkuTransacted(
	co *string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*sku.Transacted] {
	if co == nil {
		co1 := u.FormatColorOptionsOut()
		co = &co1
	}

	return sku_fmt.MakeCliFormat(
		u.config.PrintOptions,
		u.StringFormatWriterObjectId(*co),
		u.StringFormatWriterMetadatei(*co),
	)
}

func (u *Env) StringFormatWriterSkuTransactedShort() interfaces.StringFormatWriter[*sku.Transacted] {
	co := string_format_writer.ColorOptions{
		OffEntirely: true,
	}

	return sku_fmt.MakeCliFormatShort(
		u.StringFormatWriterObjectId(co),
		u.StringFormatWriterMetadatei(co),
	)
}
