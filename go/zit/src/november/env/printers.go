package env

import (
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/id_fmts"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/chrome"
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
		u.store.GetAbbrStore().Shas().Abbreviate,
	)
}

func (u *Env) StringFormatWriterKennungAligned(
	co string_format_writer.ColorOptions,
) id_fmts.Aligned {
	return id_fmts.MakeAligned(
		u.config.PrintOptions,
		u.GetStore().GetAbbrStore().GetAbbr(),
	)
}

func (u *Env) StringFormatWriterKennung(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*ids.ObjectId] {
	return id_fmts.MakeKennungCliFormat(
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

func (u *Env) SkuFmtOrganize() *sku_fmt.Organize {
	co := u.FormatColorOptionsOut()
	co.OffEntirely = true

	return sku_fmt.MakeFormatOrganize(
		u.config.PrintOptions,
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterKennungAligned(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterDescription(descriptions.CliFormatTruncationNone, co, false),
		u.StringFormatWriterEtiketten(co),
	)
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
		u.StringFormatWriterKennung(*co),
		u.StringFormatWriterMetadatei(*co),
	)
}

func (u *Env) StringFormatWriterSkuTransactedShort() interfaces.StringFormatWriter[*sku.Transacted] {
	co := string_format_writer.ColorOptions{
		OffEntirely: true,
	}

	return sku_fmt.MakeCliFormatShort(
		u.StringFormatWriterKennung(co),
		u.StringFormatWriterMetadatei(co),
	)
}

func (u *Env) PrinterSkuTransacted() interfaces.FuncIter[*sku.Transacted] {
	sw := u.StringFormatWriterSkuTransacted(nil)

	return string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		sw,
	)
}

func (u *Env) PrinterTransactedLike() interfaces.FuncIter[*sku.Transacted] {
	sw := u.StringFormatWriterSkuTransacted(nil)

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

func (u *Env) GetTime() time.Time {
	return time.Now()
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
					string_format_writer.ColorTypeTitle,
				),
			),
		)
	} else {
		return func(v string) error { return ui.Log().Print(v) }
	}
}

func (u *Env) PrinterCheckedOutFS() interfaces.FuncIter[sku.CheckedOutLike] {
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
			u.StringFormatWriterKennung(oo.ColorOptionsErr),
			u.StringFormatWriterMetadatei(oo.ColorOptionsErr),
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
			u.StringFormatWriterKennung(oo.ColorOptionsOut),
			u.StringFormatWriterMetadatei(oo.ColorOptionsOut),
		),
	)

	return func(co sku.CheckedOutLike) error {
		if co.GetState() == checked_out_state.StateError {
			return err(co)
		} else {
			return out(co)
		}
	}
}

func (u *Env) PrinterCheckedOutChrome() interfaces.FuncIter[sku.CheckedOutLike] {
	oo := u.FormatOutputOptions()

	err := string_format_writer.MakeDelim(
		"\n",
		u.Err(),
		chrome.MakeCliCheckedOutFormat(
			u.config.PrintOptions,
			u.StringFormatWriterShaLike(oo.ColorOptionsErr),
			u.StringFormatWriterKennung(oo.ColorOptionsErr),
			u.StringFormatWriterMetadatei(oo.ColorOptionsErr),
			u.StringFormatWriterTyp(oo.ColorOptionsErr),
			u.StringFormatWriterDescription(
				descriptions.CliFormatTruncation66CharEllipsis,
				oo.ColorOptionsErr,
				true,
			),
			u.StringFormatWriterEtiketten(oo.ColorOptionsErr),
		),
	)

	out := string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		chrome.MakeCliCheckedOutFormat(
			u.config.PrintOptions,
			u.StringFormatWriterShaLike(oo.ColorOptionsOut),
			u.StringFormatWriterKennung(oo.ColorOptionsOut),
			u.StringFormatWriterMetadatei(oo.ColorOptionsOut),
			u.StringFormatWriterTyp(oo.ColorOptionsOut),
			u.StringFormatWriterDescription(
				descriptions.CliFormatTruncation66CharEllipsis,
				oo.ColorOptionsOut,
				true,
			),
			u.StringFormatWriterEtiketten(oo.ColorOptionsOut),
		),
	)

	return func(co sku.CheckedOutLike) error {
		if co.GetState() == checked_out_state.StateError {
			return err(co)
		} else {
			return out(co)
		}
	}
}

func (u *Env) PrinterCheckedOutForKasten(
	k ids.RepoId,
) interfaces.FuncIter[sku.CheckedOutLike] {
	pcofs := u.PrinterCheckedOutFS()
	pcochrome := u.PrinterCheckedOutChrome()

	switch k.GetRepoIdString() {
	case "chrome":
		return pcochrome

	default:
		return pcofs
	}
}

func (u *Env) PrinterCheckedOutLike() interfaces.FuncIter[sku.CheckedOutLike] {
	pcofs := u.PrinterCheckedOutFS()
	pcochrome := u.PrinterCheckedOutChrome()

	return func(co sku.CheckedOutLike) (err error) {
		switch co.GetRepoId().GetRepoIdString() {
		case "chrome":
			return pcochrome(co)

		default:
			return pcofs(co)
		}
	}
}

func (u *Env) PrinterMatching() sku.IterMatching {
	pt := u.PrinterSkuTransacted()
	pco := u.PrinterCheckedOutLike()

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
					checked_out_state.StateRecognized,
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = co.GetSku().SetFromSkuLike(sk); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = pco(co); err != nil {
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
