package umwelt

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_fmt"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
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
		u.storeUtil.GetAbbrStore().Shas().Abbreviate,
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

func (u *Umwelt) PrinterHeader() schnittstellen.FuncIter[string] {
	return string_format_writer.MakeDelim[string](
		"\n",
		u.Out(),
		string_format_writer.MakeIndentedHeader(u.FormatColorOptionsOut()),
	)
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