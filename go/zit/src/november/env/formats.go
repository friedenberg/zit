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
	"code.linenisgreat.com/zit/go/zit/src/lima/store_browser"
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
	truncate string_format_writer.CliFormatTruncation,
	co string_format_writer.ColorOptions,
	quote bool,
) interfaces.StringFormatWriter[*descriptions.Description] {
	return descriptions.MakeCliFormat(truncate, co, quote)
}

func (u *Env) StringFormatWriterFields(
	truncate string_format_writer.CliFormatTruncation,
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[[]string_format_writer.Field] {
	return string_format_writer.MakeCliFormatFields(truncate, co)
}

func (u *Env) StringFormatWriterEtiketten(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*ids.Tag] {
	return id_fmts.MakeEtikettenCliFormat()
}

func (u *Env) StringFormatWriterMetadatei(
	co string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
) interfaces.StringFormatWriter[*object_metadata.Metadata] {
	return sku_fmt.MakeCliMetadateiFormat(
		u.config.PrintOptions,
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterFields(truncation, co),
		u.StringFormatWriterEtiketten(co),
	)
}

func (u *Env) StringFormatWriterSku(
	co string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
) *sku_fmt.Organize {
	return sku_fmt.MakeFormatOrganize(
		co,
		u.config.PrintOptions,
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterObjectIdAligned(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterEtiketten(co),
		u.StringFormatWriterFields(truncation, co),
		u.StringFormatWriterMetadatei(
			co,
			truncation,
		),
	)
}

func (u *Env) SkuFmtOrganize(repoId ids.RepoId) sku_fmt.ExternalLike {
	co := u.FormatColorOptionsOut()
	co.OffEntirely = true

	f := u.StringFormatWriterSku(
		co,
		string_format_writer.CliFormatTruncationNone,
	)

	kid := repoId.GetRepoIdString()
	es, ok := u.externalStores[kid]

	if !ok {
		return sku_fmt.ExternalLike{
			ReaderExternalLike: f,
			WriterExternalLike: f,
		}
	}

	return es.GetExternalStoreOrganizeFormat(f)
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
		u.StringFormatWriterMetadatei(*co, truncate),
	)
}

func (u *Env) StringFormatWriterSkuTransactedShort() interfaces.StringFormatWriter[*sku.Transacted] {
	co := string_format_writer.ColorOptions{
		OffEntirely: true,
	}

	return sku_fmt.MakeCliFormatShort(
		u.StringFormatWriterObjectId(co),
		u.StringFormatWriterMetadatei(
			co,
			string_format_writer.CliFormatTruncation66CharEllipsis,
		),
	)
}

func (u *Env) StringFormatWriterStoreBrowserCheckedOut() interfaces.StringFormatWriter[sku.CheckedOutLike] {
	return store_browser.MakeCliCheckedOutFormat(
		u.config.PrintOptions,
		store_browser.MakeFormatOrganize(
			u.StringFormatWriterSku(
				u.FormatColorOptionsOut(),
				string_format_writer.CliFormatTruncation66CharEllipsis,
			),
		),
	)
}
