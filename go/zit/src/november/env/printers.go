package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/id_fmts"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
)

func (u *Env) PrinterSkuTransacted() interfaces.FuncIter[*sku.Transacted] {
	sw := u.StringFormatWriterSkuTransacted(
		nil,
		string_format_writer.CliFormatTruncation66CharEllipsis,
	)

	return string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		sw,
	)
}

func (u *Env) PrinterTransactedLike() interfaces.FuncIter[*sku.Transacted] {
	sw := u.StringFormatWriterSkuTransacted(
		nil,
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

func (u *Env) PrinterFileNotRecognized() interfaces.FuncIter[*fd.FD] {
	p := id_fmts.MakeFileNotRecognizedStringWriterFormat(
		id_fmts.MakeFDCliFormat(
			u.FormatColorOptionsOut(),
			u.fs_home.MakeRelativePathStringFormatWriter(),
		),
		u.StringFormatWriterShaLike(u.FormatColorOptionsOut()),
	)

	return string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		p,
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

func (u *Env) PrinterCheckedOutFS() interfaces.FuncIter[*sku.CheckedOut] {
	oo := u.FormatOutputOptions()

	err := string_format_writer.MakeDelim(
		"\n",
		u.Err(),
		store_fs.MakeCliCheckedOutFormat(
			u.config.PrintOptions,
			u.StringFormatWriterShaLike(oo.ColorOptionsErr),
			id_fmts.MakeFDCliFormat(
				oo.ColorOptionsErr,
				u.fs_home.MakeRelativePathStringFormatWriter(),
			),
			u.StringFormatWriterObjectId(oo.ColorOptionsErr),
			u.StringFormatWriterMetadata(
				oo.ColorOptionsErr,
				string_format_writer.CliFormatTruncation66CharEllipsis,
			),
			u.GetStore().GetCwdFiles(),
		),
	)

	out := string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		store_fs.MakeCliCheckedOutFormat(
			u.config.PrintOptions,
			u.StringFormatWriterShaLike(oo.ColorOptionsOut),
			id_fmts.MakeFDCliFormat(
				oo.ColorOptionsOut,
				u.fs_home.MakeRelativePathStringFormatWriter(),
			),
			u.StringFormatWriterObjectId(oo.ColorOptionsOut),
			u.StringFormatWriterMetadata(
				oo.ColorOptionsErr,
				string_format_writer.CliFormatTruncation66CharEllipsis,
			),
			u.GetStore().GetCwdFiles(),
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

func (u *Env) PrinterCheckedOutBrowser() interfaces.FuncIter[sku.CheckedOutLike] {
	sfw := u.StringFormatWriterStoreBrowserCheckedOut()

	err := string_format_writer.MakeDelim(
		"\n",
		u.Err(),
		sfw,
	)

	out := string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		sfw,
	)

	return func(co sku.CheckedOutLike) error {
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
	return u.PrinterCheckedOutFS()
}

func (u *Env) PrinterMatching() sku.IterMatching {
	pt := u.PrinterSkuTransacted()
	pco := u.PrinterCheckedOutFS()

	return func(
		mt sku.UnsureMatchType,
		sk *sku.Transacted,
		existing sku.CheckedOutLikeMutableSet,
	) (err error) {
		if err = pt(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = existing.Each(
			func(co sku.CheckedOutLike) (err error) {
				if err = co.SetState(
					checked_out_state.Recognized,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				sku.TransactedResetter.ResetWith(co.GetSku(), sk)

				if err = pco(co.(*sku.CheckedOut)); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
