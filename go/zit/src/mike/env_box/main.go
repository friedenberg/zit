package env_box

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
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

	PrinterTransacted() interfaces.FuncIter[*sku.Transacted]
	PrinterCheckedOut(
		headerWriter string_format_writer.HeaderWriter[*sku.CheckedOut],
	) interfaces.FuncIter[*sku.CheckedOut]

	GetUIStorePrinters() sku.UIStorePrinters
}

// TODO make this work even if storeFS and abbr are nil
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

func (env *env) GetAbbr() (abbr ids.Abbr) {
	if env.abbr != nil {
		abbr = env.abbr.GetAbbr()
	}

	return
}

func (env *env) StringFormatWriterSkuBoxTransacted(
	printOptions options_print.V0,
	colorOptions string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
) *box_format.BoxTransacted {
	var headerWriter string_format_writer.HeaderWriter[*sku.Transacted]

	if printOptions.PrintTime && !printOptions.PrintTai {
		headerWriter = box_format.TransactedHeaderUserTai{}
	}

	return box_format.MakeBoxTransacted(
		colorOptions,
		printOptions,
		env.StringFormatWriterFields(truncation, colorOptions),
		env.GetAbbr(),
		env.storeFS,
		env,
		headerWriter,
	)
}

func (env *env) StringFormatWriterSkuBoxCheckedOut(
	printOptions options_print.V0,
	colorOptions string_format_writer.ColorOptions,
	truncation string_format_writer.CliFormatTruncation,
	headerWriter string_format_writer.HeaderWriter[*sku.CheckedOut],
) *box_format.BoxCheckedOut {
	return box_format.MakeBoxCheckedOut(
		colorOptions,
		printOptions,
		env.StringFormatWriterFields(truncation, colorOptions),
		env.GetAbbr(),
		env.storeFS,
		env,
		headerWriter,
	)
}

func (env *env) SkuFormatBoxTransactedNoColor() *box_format.BoxTransacted {
	colorOptions := env.FormatColorOptionsOut()
	colorOptions.OffEntirely = true
	options := env.GetCLIConfig().PrintOptions.WithPrintShas(false)
	options.PrintTime = false
	options.PrintShas = false
	options.DescriptionInBox = false

	return env.StringFormatWriterSkuBoxTransacted(
		options,
		colorOptions,
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

func (env *env) PrinterTransacted() interfaces.FuncIter[*sku.Transacted] {
	printOptions := env.GetCLIConfig().PrintOptions.
		WithPrintShas(true).
		WithExcludeFields(true)

	stringFormatWriter := env.StringFormatWriterSkuBoxTransacted(
		printOptions,
		env.FormatColorOptionsOut(),
		string_format_writer.CliFormatTruncation66CharEllipsis,
	)

	return string_format_writer.MakeDelim(
		"\n",
		env.GetUIFile(),
		string_format_writer.MakeFunc(
			func(
				writer interfaces.WriterAndStringWriter,
				object *sku.Transacted,
			) (n int64, err error) {
				return stringFormatWriter.EncodeStringTo(object, writer)
			},
		),
	)
}

func (env *env) PrinterCheckedOut(
	headerWriter string_format_writer.HeaderWriter[*sku.CheckedOut],
) interfaces.FuncIter[*sku.CheckedOut] {
	oo := env.FormatOutputOptions()
	po := env.GetCLIConfig().PrintOptions.
		WithPrintShas(true)

	out := string_format_writer.MakeDelim(
		"\n",
		env.GetUIFile(),
		env.StringFormatWriterSkuBoxCheckedOut(
			po,
			oo.ColorOptionsErr,
			string_format_writer.CliFormatTruncation66CharEllipsis,
			box_format.CheckedOutHeaderState{},
		),
	)

	return out
}

func (env *env) GetUIStorePrinters() sku.UIStorePrinters {
	printerTransacted := env.PrinterTransacted()

	return sku.UIStorePrinters{
		TransactedNew:     printerTransacted,
		TransactedUpdated: printerTransacted,
		TransactedUnchanged: func(sk *sku.Transacted) (err error) {
			if !env.GetCLIConfig().PrintOptions.PrintUnchanged {
				return
			}

			return printerTransacted(sk)
		},
		CheckedOutCheckedOut: env.PrinterCheckedOut(
			box_format.CheckedOutHeaderState{},
		),
	}
}
