package repo_local

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/id_fmts"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
)

func (repo *Repo) PrinterTransacted() interfaces.FuncIter[*sku.Transacted] {
	po := repo.config.PrintOptions.
		WithPrintShas(true).
		WithExcludeFields(true)

	sw := repo.StringFormatWriterSkuBoxTransacted(
		po,
		repo.FormatColorOptionsOut(),
		string_format_writer.CliFormatTruncation66CharEllipsis,
	)

	return string_format_writer.MakeDelim(
		"\n",
		repo.GetUIFile(),
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return sw.WriteStringFormat(w, o)
			},
		),
	)
}

// TODO migrate to StringFormatWriterSkuBoxCheckedOut
func (repo *Repo) PrinterTransactedDeleted() interfaces.FuncIter[*sku.CheckedOut] {
	po := repo.config.PrintOptions.
		WithPrintShas(true).
		WithPrintTime(false)

	sw := repo.StringFormatWriterSkuBoxCheckedOut(
		po,
		repo.FormatColorOptionsOut(),
		string_format_writer.CliFormatTruncation66CharEllipsis,
		box_format.CheckedOutHeaderDeleted{
			ConfigDryRunReader: repo.GetConfig(),
		},
	)

	return string_format_writer.MakeDelim(
		"\n",
		repo.GetUIFile(),
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.CheckedOut) (n int64, err error) {
				return sw.WriteStringFormat(w, o)
			},
		),
	)
}

// TODO make generic external version
func (u *Repo) PrinterFDDeleted() interfaces.FuncIter[*fd.FD] {
	p := id_fmts.MakeFDDeletedStringWriterFormat(
		u.GetConfig().DryRun,
		id_fmts.MakeFDCliFormat(
			u.FormatColorOptionsOut(),
			u.layout.MakeRelativePathStringFormatWriter(),
		),
	)

	return string_format_writer.MakeDelim(
		"\n",
		u.GetUIFile(),
		p,
	)
}

func (u *Repo) PrinterHeader() interfaces.FuncIter[string] {
	if u.config.PrintOptions.PrintFlush {
		return string_format_writer.MakeDelim(
			"\n",
			u.GetErrFile(),
			string_format_writer.MakeDefaultDatePrefixFormatWriter(
				u,
				string_format_writer.MakeColor(
					u.FormatColorOptionsOut(),
					string_format_writer.MakeString[string](),
					string_format_writer.ColorTypeHeading,
				),
			),
		)
	} else {
		return func(v string) error { return ui.Log().Print(v) }
	}
}

func (u *Repo) PrinterCheckedOut(
	headerWriter string_format_writer.HeaderWriter[*sku.CheckedOut],
) interfaces.FuncIter[*sku.CheckedOut] {
	oo := u.FormatOutputOptions()
	po := u.config.PrintOptions.
		WithPrintShas(true)

	out := string_format_writer.MakeDelim(
		"\n",
		u.GetUIFile(),
		u.StringFormatWriterSkuBoxCheckedOut(
			po,
			oo.ColorOptionsErr,
			string_format_writer.CliFormatTruncation66CharEllipsis,
			box_format.CheckedOutHeaderState{},
		),
	)

	return out
}

func (u *Repo) MakePrinterBoxArchive(
	out interfaces.WriterAndStringWriter,
	includeTai bool,
) interfaces.FuncIter[*sku.Transacted] {
	boxFormat := u.MakeBoxArchive(includeTai)

	return string_format_writer.MakeDelim(
		"\n",
		out,
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return boxFormat.WriteStringFormat(w, o)
			},
		),
	)
}
