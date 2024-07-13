package umwelt

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
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/kennung_fmt"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/chrome"
)

func (u *Umwelt) FormatOutputOptions() (o string_format_writer.OutputOptions) {
	o.ColorOptionsOut = u.FormatColorOptionsOut()
	o.ColorOptionsErr = u.FormatColorOptionsErr()
	return
}

func (u *Umwelt) FormatColorOptionsOut() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.outIsTty || !u.konfig.PrintOptions.PrintColors
	return
}

func (u *Umwelt) FormatColorOptionsErr() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.errIsTty || !u.konfig.PrintOptions.PrintColors
	return
}

func (u *Umwelt) StringFormatWriterShaLike(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[interfaces.ShaLike] {
	return kennung_fmt.MakeShaCliFormat(
		u.konfig.PrintOptions,
		co,
		u.store.GetAbbrStore().Shas().Abbreviate,
	)
}

func (u *Umwelt) StringFormatWriterKennungAligned(
	co string_format_writer.ColorOptions,
) kennung_fmt.Aligned {
	return kennung_fmt.MakeAligned(
		u.konfig.PrintOptions,
		u.GetStore().GetAbbrStore().GetAbbr(),
	)
}

func (u *Umwelt) StringFormatWriterKennung(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*ids.ObjectId] {
	return kennung_fmt.MakeKennungCliFormat(
		u.konfig.PrintOptions,
		co,
		u.GetStore().GetAbbrStore().GetAbbr(),
	)
}

func (u *Umwelt) StringFormatWriterTyp(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*ids.Type] {
	return kennung_fmt.MakeTypCliFormat(co)
}

func (u *Umwelt) StringFormatWriterBezeichnung(
	truncate descriptions.CliFormatTruncation,
	co string_format_writer.ColorOptions,
	quote bool,
) interfaces.StringFormatWriter[*descriptions.Description] {
	return descriptions.MakeCliFormat(truncate, co, quote)
}

func (u *Umwelt) StringFormatWriterEtiketten(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*ids.Tag] {
	return kennung_fmt.MakeEtikettenCliFormat()
}

func (u *Umwelt) StringFormatWriterMetadatei(
	co string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*metadatei.Metadatei] {
	return sku_fmt.MakeCliMetadateiFormat(
		u.konfig.PrintOptions,
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterBezeichnung(
			descriptions.CliFormatTruncation66CharEllipsis,
			co,
			true,
		),
		u.StringFormatWriterEtiketten(co),
	)
}

func (u *Umwelt) SkuFmtOrganize() *sku_fmt.Organize {
	co := u.FormatColorOptionsOut()
	co.OffEntirely = true

	return sku_fmt.MakeFormatOrganize(
		u.konfig.PrintOptions,
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterKennungAligned(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterBezeichnung(descriptions.CliFormatTruncationNone, co, false),
		u.StringFormatWriterEtiketten(co),
	)
}

func (u *Umwelt) StringFormatWriterSkuTransacted(
	co *string_format_writer.ColorOptions,
) interfaces.StringFormatWriter[*sku.Transacted] {
	if co == nil {
		co1 := u.FormatColorOptionsOut()
		co = &co1
	}

	return sku_fmt.MakeCliFormat(
		u.konfig.PrintOptions,
		u.StringFormatWriterKennung(*co),
		u.StringFormatWriterMetadatei(*co),
	)
}

func (u *Umwelt) StringFormatWriterSkuTransactedShort() interfaces.StringFormatWriter[*sku.Transacted] {
	co := string_format_writer.ColorOptions{
		OffEntirely: true,
	}

	return sku_fmt.MakeCliFormatShort(
		u.StringFormatWriterKennung(co),
		u.StringFormatWriterMetadatei(co),
	)
}

func (u *Umwelt) PrinterSkuTransacted() interfaces.FuncIter[*sku.Transacted] {
	sw := u.StringFormatWriterSkuTransacted(nil)

	return string_format_writer.MakeDelim(
		"\n",
		u.Out(),
		sw,
	)
}

func (u *Umwelt) PrinterTransactedLike() interfaces.FuncIter[*sku.Transacted] {
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

func (u *Umwelt) PrinterFileNotRecognized() interfaces.FuncIter[*fd.FD] {
	p := kennung_fmt.MakeFileNotRecognizedStringWriterFormat(
		kennung_fmt.MakeFDCliFormat(
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
func (u *Umwelt) PrinterFDDeleted() interfaces.FuncIter[*fd.FD] {
	p := kennung_fmt.MakeFDDeletedStringWriterFormat(
		u.GetKonfig().DryRun,
		kennung_fmt.MakeFDCliFormat(
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

func (u *Umwelt) GetTime() time.Time {
	return time.Now()
}

func (u *Umwelt) PrinterHeader() interfaces.FuncIter[string] {
	if u.konfig.PrintOptions.PrintFlush {
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

func (u *Umwelt) PrinterCheckedOutFS() interfaces.FuncIter[sku.CheckedOutLike] {
	oo := u.FormatOutputOptions()

	err := string_format_writer.MakeDelim(
		"\n",
		u.Err(),
		store_fs.MakeCliCheckedOutFormat(
			u.konfig.PrintOptions,
			u.StringFormatWriterShaLike(oo.ColorOptionsErr),
			kennung_fmt.MakeFDCliFormat(
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
			u.konfig.PrintOptions,
			u.StringFormatWriterShaLike(oo.ColorOptionsOut),
			kennung_fmt.MakeFDCliFormat(
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

func (u *Umwelt) PrinterCheckedOutChrome() interfaces.FuncIter[sku.CheckedOutLike] {
	oo := u.FormatOutputOptions()

	err := string_format_writer.MakeDelim(
		"\n",
		u.Err(),
		chrome.MakeCliCheckedOutFormat(
			u.konfig.PrintOptions,
			u.StringFormatWriterShaLike(oo.ColorOptionsErr),
			u.StringFormatWriterKennung(oo.ColorOptionsErr),
			u.StringFormatWriterMetadatei(oo.ColorOptionsErr),
			u.StringFormatWriterTyp(oo.ColorOptionsErr),
			u.StringFormatWriterBezeichnung(
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
			u.konfig.PrintOptions,
			u.StringFormatWriterShaLike(oo.ColorOptionsOut),
			u.StringFormatWriterKennung(oo.ColorOptionsOut),
			u.StringFormatWriterMetadatei(oo.ColorOptionsOut),
			u.StringFormatWriterTyp(oo.ColorOptionsOut),
			u.StringFormatWriterBezeichnung(
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

func (u *Umwelt) PrinterCheckedOutForKasten(
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

func (u *Umwelt) PrinterCheckedOutLike() interfaces.FuncIter[sku.CheckedOutLike] {
	pcofs := u.PrinterCheckedOutFS()
	pcochrome := u.PrinterCheckedOutChrome()

	return func(co sku.CheckedOutLike) (err error) {
		switch co.GetKasten().GetRepoIdString() {
		case "chrome":
			return pcochrome(co)

		default:
			return pcofs(co)
		}
	}
}

func (u *Umwelt) PrinterMatching() sku.IterMatching {
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
