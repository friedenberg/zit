package env_box

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
)

type Env interface {
	StringFormatWriterSkuBoxTransacted(
		po options_print.V0,
		co string_format_writer.ColorOptions,
		truncation string_format_writer.CliFormatTruncation,
	) *box_format.BoxTransacted

	StringFormatWriterSkuBoxCheckedOut(
		po options_print.V0,
		co string_format_writer.ColorOptions,
		truncation string_format_writer.CliFormatTruncation,
		headerWriter string_format_writer.HeaderWriter[*sku.CheckedOut],
	) *box_format.BoxCheckedOut

	SkuFormatBoxTransactedNoColor() *box_format.BoxTransacted
	SkuFormatBoxCheckedOutNoColor() *box_format.BoxCheckedOut
}

func Make(
	envRepo env_repo.Env,
	storeFS *store_fs.Store,
	abbr sku.AbbrStore,
) Env {
	return &env{
		Env:     envRepo,
		storeFS: storeFS,
		abbr:    abbr,
	}
}

type env struct {
	env_repo.Env
	storeFS       *store_fs.Store
	abbr          sku.AbbrStore
	object_format object_inventory_format.Format
	options       object_inventory_format.Options
}

func (s env) FormatForVersion(sv interfaces.StoreVersion) sku.ListFormat {
	v := sv.GetInt()

	switch {
	case v <= 6:
		return inventory_list_blobs.MakeV0(
			s.object_format,
			s.options,
		)

	default:
		return inventory_list_blobs.V1{
			Box: box_format.MakeBoxTransactedArchive(
				s,
				options_print.V0{}.WithPrintTai(true),
			),
		}
	}
}

func (u *env) StringFormatWriterSkuBoxTransacted(
	po options_print.V0,
	co string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
) *box_format.BoxTransacted {
	var headerWriter string_format_writer.HeaderWriter[*sku.Transacted]

	if po.PrintTime && !po.PrintTai {
		headerWriter = box_format.TransactedHeaderUserTai{}
	}

	return box_format.MakeBoxTransacted(
		co,
		po,
		u.StringFormatWriterFields(truncation, co),
		u.abbr.GetAbbr(),
		u.storeFS,
		u,
		headerWriter,
	)
}

func (u *env) StringFormatWriterSkuBoxCheckedOut(
	po options_print.V0,
	co string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
	headerWriter string_format_writer.HeaderWriter[*sku.CheckedOut],
) *box_format.BoxCheckedOut {
	return box_format.MakeBoxCheckedOut(
		co,
		po,
		u.StringFormatWriterFields(truncation, co),
		u.abbr.GetAbbr(),
		u.storeFS,
		u,
		headerWriter,
	)
}

func (u *env) SkuFormatBoxTransactedNoColor() *box_format.BoxTransacted {
	co := u.FormatColorOptionsOut()
	co.OffEntirely = true
	options := u.GetCLIConfig().PrintOptions.WithPrintShas(false)
	options.PrintTime = false
	options.PrintShas = false
	options.DescriptionInBox = false

	return u.StringFormatWriterSkuBoxTransacted(
		options,
		co,
		string_format_writer.CliFormatTruncationNone,
	)
}

func (u *env) SkuFormatBoxCheckedOutNoColor() *box_format.BoxCheckedOut {
	co := u.FormatColorOptionsOut()
	co.OffEntirely = true
	options := u.GetCLIConfig().PrintOptions.WithPrintShas(false)
	options.PrintTime = false
	options.PrintShas = false
	options.DescriptionInBox = false

	return u.StringFormatWriterSkuBoxCheckedOut(
		options,
		co,
		string_format_writer.CliFormatTruncationNone,
		nil,
	)
}
