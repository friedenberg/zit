package umwelt

import (
	"time"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/catgut"
	"code.linenisgreat.com/zit/src/charlie/string_format_writer"
	"code.linenisgreat.com/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/src/echo/bezeichnung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/kennung_fmt"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
)

func (u *Umwelt) FormatColorOptionsOut() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.outIsTty
	return
}

func (u *Umwelt) FormatColorOptionsErr() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.errIsTty
	return
}

func (u *Umwelt) StringFormatWriterShaLike(
	co string_format_writer.ColorOptions,
) schnittstellen.StringFormatWriter[schnittstellen.ShaLike] {
	return kennung_fmt.MakeShaCliFormat(
		u.konfig.PrintOptions,
		co,
		u.store.GetAbbrStore().Shas().Abbreviate,
	)
}

func (u *Umwelt) StringFormatWriterKennung(
	co string_format_writer.ColorOptions,
) schnittstellen.StringFormatWriter[*kennung.Kennung2] {
	return kennung_fmt.MakeKennungCliFormat(
		u.konfig.PrintOptions,
		co,
		u.MakeKennungExpanders(),
	)
}

func (u *Umwelt) StringFormatWriterTyp(
	co string_format_writer.ColorOptions,
) schnittstellen.StringFormatWriter[*kennung.Typ] {
	return kennung_fmt.MakeTypCliFormat(co)
}

func (u *Umwelt) StringFormatWriterBezeichnung(
	truncate bezeichnung.CliFormatTruncation,
	co string_format_writer.ColorOptions,
) schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung] {
	return bezeichnung.MakeCliFormat2(truncate, co)
}

func (u *Umwelt) StringFormatWriterEtiketten(
	co string_format_writer.ColorOptions,
) schnittstellen.StringFormatWriter[*kennung.Etikett] {
	return kennung_fmt.MakeEtikettenCliFormat()
}

func (u *Umwelt) SkuFmtNewOrganize() *sku_fmt.OrganizeNew {
	co := u.FormatColorOptionsOut()
	co.OffEntirely = true

	return sku_fmt.MakeOrganizeNewFormat(
		u.konfig.PrintOptions,
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterKennung(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterBezeichnung(bezeichnung.CliFormatTruncationNone, co),
		u.StringFormatWriterEtiketten(co),
	)
}

func (u *Umwelt) SkuFormatOldOrganize() *sku_fmt.Organize {
	return sku_fmt.MakeOrganizeFormat(
		u.MakeKennungExpanders(),
		u.konfig.PrintOptions,
	)
}

func (u *Umwelt) StringFormatWriterSkuLikePtrForOrganize() catgut.StringFormatReadWriter[*sku.Transacted] {
	if !u.Konfig().NewOrganize {
		return u.SkuFormatOldOrganize()
	}

	return u.SkuFmtNewOrganize()
}

func (u *Umwelt) StringFormatWriterSkuTransacted(
	co *string_format_writer.ColorOptions,
) schnittstellen.StringFormatWriter[*sku.Transacted] {
	if co == nil {
		co1 := u.FormatColorOptionsOut()
		co = &co1
	}

	return sku_fmt.MakeCliFormat(
		u.konfig.PrintOptions,
		u.StringFormatWriterShaLike(*co),
		u.StringFormatWriterKennung(*co),
		u.StringFormatWriterTyp(*co),
		u.StringFormatWriterBezeichnung(
			bezeichnung.CliFormatTruncation66CharEllipsis,
			*co,
		),
		u.StringFormatWriterEtiketten(*co),
	)
}

func (u *Umwelt) StringFormatWriterSkuTransactedShort() schnittstellen.StringFormatWriter[*sku.Transacted] {
	co := string_format_writer.ColorOptions{
		OffEntirely: true,
	}

	return sku_fmt.MakeCliFormatShort(
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterKennung(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterBezeichnung(
			bezeichnung.CliFormatTruncation66CharEllipsis,
			co,
		),
		u.StringFormatWriterEtiketten(co),
	)
}

func (u *Umwelt) PrinterSkuTransacted() schnittstellen.FuncIter[*sku.Transacted] {
	sw := u.StringFormatWriterSkuTransacted(nil)

	return string_format_writer.MakeDelim[*sku.Transacted](
		"\n",
		u.Out(),
		sw,
	)
}

func (u *Umwelt) PrinterTransactedLike() schnittstellen.FuncIter[*sku.Transacted] {
	sw := u.StringFormatWriterSkuTransacted(nil)

	return string_format_writer.MakeDelim[*sku.Transacted](
		"\n",
		u.Out(),
		string_format_writer.MakeFunc[*sku.Transacted](
			func(w schnittstellen.WriterAndStringWriter, o *sku.Transacted) (n int64, err error) {
				return sw.WriteStringFormat(w, o)
			},
		),
	)
}

func (u *Umwelt) PrinterFileNotRecognized() schnittstellen.FuncIter[*fd.FD] {
	p := kennung_fmt.MakeFileNotRecognizedStringWriterFormat(
		kennung_fmt.MakeFDCliFormat(
			u.FormatColorOptionsOut(),
			u.standort.MakeRelativePathStringFormatWriter(),
		),
		u.StringFormatWriterShaLike(u.FormatColorOptionsOut()),
	)

	return string_format_writer.MakeDelim[*fd.FD](
		"\n",
		u.Out(),
		p,
	)
}

func (u *Umwelt) PrinterFDDeleted() schnittstellen.FuncIter[*fd.FD] {
	p := kennung_fmt.MakeFDDeletedStringWriterFormat(
		u.Konfig().DryRun,
		kennung_fmt.MakeFDCliFormat(
			u.FormatColorOptionsOut(),
			u.standort.MakeRelativePathStringFormatWriter(),
		),
	)

	return string_format_writer.MakeDelim[*fd.FD](
		"\n",
		u.Out(),
		p,
	)
}

func (u *Umwelt) GetTime() time.Time {
	return time.Now()
}

func (u *Umwelt) PrinterHeader() schnittstellen.FuncIter[string] {
	if u.konfig.PrintOptions.PrintFlush {
		return string_format_writer.MakeDelim[string](
			"\n",
			u.Err(),
			string_format_writer.MakeDefaultDatePrefixFormatWriter(
				u,
				string_format_writer.MakeColor[string](
					u.FormatColorOptionsOut(),
					string_format_writer.MakeString[string](),
					string_format_writer.ColorTypeTitle,
				),
			),
		)
	} else {
		return func(string) error { return nil }
	}
}

func (u *Umwelt) PrinterCheckedOut() schnittstellen.FuncIter[*sku.CheckedOut] {
	coOut := u.FormatColorOptionsOut()
	coErr := u.FormatColorOptionsErr()

	err := string_format_writer.MakeDelim[*sku.CheckedOut](
		"\n",
		u.Err(),
		sku_fmt.MakeCliCheckedOutFormat(
			u.konfig.PrintOptions,
			u.StringFormatWriterShaLike(coErr),
			kennung_fmt.MakeFDCliFormat(
				coErr,
				u.standort.MakeRelativePathStringFormatWriter(),
			),
			u.StringFormatWriterKennung(coErr),
			u.StringFormatWriterTyp(coErr),
			u.StringFormatWriterBezeichnung(
				bezeichnung.CliFormatTruncation66CharEllipsis,
				coErr,
			),
			u.StringFormatWriterEtiketten(coErr),
		),
	)

	out := string_format_writer.MakeDelim[*sku.CheckedOut](
		"\n",
		u.Out(),
		sku_fmt.MakeCliCheckedOutFormat(
			u.konfig.PrintOptions,
			u.StringFormatWriterShaLike(coOut),
			kennung_fmt.MakeFDCliFormat(
				coOut,
				u.standort.MakeRelativePathStringFormatWriter(),
			),
			u.StringFormatWriterKennung(coOut),
			u.StringFormatWriterTyp(coOut),
			u.StringFormatWriterBezeichnung(
				bezeichnung.CliFormatTruncation66CharEllipsis,
				coOut,
			),
			u.StringFormatWriterEtiketten(coOut),
		),
	)

	return func(co *sku.CheckedOut) error {
		if co.State == checked_out_state.StateError {
			return err(co)
		} else {
			return out(co)
		}
	}
}
