package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/id_fmts"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (u *Env) PrinterTransacted() interfaces.FuncIter[*sku.Transacted] {
	po := u.config.PrintOptions.
		WithPrintShas(true).
		WithExcludeFields(true)

	sw := u.StringFormatWriterSkuBox(
		po,
		u.FormatColorOptionsOut(),
		string_format_writer.CliFormatTruncation66CharEllipsis,
	)

	return string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return sw.WriteStringFormat(w, o)
			},
		),
	)
}

func (u *Env) PrinterTransactedDeleted() interfaces.FuncIter[*sku.Transacted] {
	po := u.config.PrintOptions.
		WithPrintShas(true).
		WithPrintTime(false).
		WithPrintState(true)

	sw := u.StringFormatWriterSkuBox(
		po,
		u.FormatColorOptionsOut(),
		string_format_writer.CliFormatTruncation66CharEllipsis,
	)

	return string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		string_format_writer.MakeFunc(
			func(w interfaces.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return sw.WriteStringFormat(w, o)
			},
		),
	)
}

// TODO make generic external version
func (u *Env) PrinterFDDeleted() interfaces.FuncIter[*fd.FD] {
	p := id_fmts.MakeFDDeletedStringWriterFormat(
		u.GetConfig().DryRun,
		id_fmts.MakeFDCliFormat(
			u.FormatColorOptionsOut(),
			u.fs_home.MakeRelativePathStringFormatWriter(),
		),
	)

	return string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		p,
	)
}

func (u *Env) PrinterHeader() interfaces.FuncIter[string] {
	if u.config.PrintOptions.PrintFlush {
		return string_format_writer.MakeDelim(
			"\n",
			u.Err(),
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

func (u *Env) PrinterCheckedOut() interfaces.FuncIter[*sku.CheckedOut] {
	oo := u.FormatOutputOptions()
	po := u.config.PrintOptions.
		WithPrintShas(true)

	err := string_format_writer.MakeDelim(
		"\n",
		u.Err(),
		u.StringFormatWriterSkuBox(
			po,
			oo.ColorOptionsErr,
			string_format_writer.CliFormatTruncation66CharEllipsis,
		),
	)

	out := string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		u.StringFormatWriterSkuBox(
			po,
			oo.ColorOptionsErr,
			string_format_writer.CliFormatTruncation66CharEllipsis,
		),
	)

	return func(co *sku.CheckedOut) error {
		if co.GetState() == checked_out_state.Error {
			return err(co)
		} else {
			return out(co)
		}
	}
}

func (u *Env) PrinterCheckedOutForKasten(
	k ids.RepoId,
) interfaces.FuncIter[*sku.CheckedOut] {
	return u.PrinterCheckedOut()
}
