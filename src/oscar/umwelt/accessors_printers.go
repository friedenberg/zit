package umwelt

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/kennung_fmt"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/sku_fmt"
)

func (u *Umwelt) FormatColorOptions() (o string_format_writer.ColorOptions) {
	o.OffEntirely = !u.outIsTty
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
) schnittstellen.StringFormatWriter[kennung.Kennung2] {
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
	co string_format_writer.ColorOptions,
) schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung] {
	return bezeichnung.MakeCliFormat2(co)
}

func (u *Umwelt) StringFormatWriterEtiketten(
	co string_format_writer.ColorOptions,
) schnittstellen.StringFormatWriter[kennung.EtikettSet] {
	return kennung_fmt.MakeEtikettenCliFormat()
}

func (u *Umwelt) StringFormatWriterSkuLikePtr() schnittstellen.StringFormatWriter[*sku.Transacted] {
	co := u.FormatColorOptions()

	return sku_fmt.MakeCliFormat(
		u.konfig.PrintOptions,
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterKennung(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterBezeichnung(co),
		u.StringFormatWriterEtiketten(co),
	)
}

func (u *Umwelt) StringFormatWriterSkuLikePtrShort() schnittstellen.StringFormatWriter[*sku.Transacted] {
	co := string_format_writer.ColorOptions{
		OffEntirely: true,
	}

	return sku_fmt.MakeCliFormatShort(
		u.StringFormatWriterShaLike(co),
		u.StringFormatWriterKennung(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterBezeichnung(co),
		u.StringFormatWriterEtiketten(co),
	)
}

func (u *Umwelt) PrinterSkuLikePtr() schnittstellen.FuncIter[*sku.Transacted] {
	sw := u.StringFormatWriterSkuLikePtr()

	return string_format_writer.MakeDelim[*sku.Transacted](
		"\n",
		u.Out(),
		sw,
	)
}

func (u *Umwelt) PrinterTransactedLike() schnittstellen.FuncIter[*sku.Transacted] {
	sw := u.StringFormatWriterSkuLikePtr()

	return string_format_writer.MakeDelim[*sku.Transacted](
		"\n",
		u.Out(),
		string_format_writer.MakeFunc[*sku.Transacted](
			func(w io.StringWriter, o *sku.Transacted) (n int64, err error) {
				return sw.WriteStringFormat(w, o)
			},
		),
	)
}

func (u *Umwelt) PrinterFileNotRecognized() schnittstellen.FuncIter[*fd.FD] {
	p := kennung_fmt.MakeFileNotRecognizedStringWriterFormat(
		kennung_fmt.MakeFDCliFormat(
			u.FormatColorOptions(),
			u.standort.MakeRelativePathStringFormatWriter(),
		),
		u.StringFormatWriterShaLike(u.FormatColorOptions()),
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
			u.FormatColorOptions(),
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
		string_format_writer.MakeIndentedHeader(u.FormatColorOptions()),
	)
}

func (u *Umwelt) PrinterCheckedOutLike() schnittstellen.FuncIter[*sku.CheckedOut] {
	co := u.FormatColorOptions()

	p := sku.MakeCliFormat(
		sku.CliOptions{},
		u.StringFormatWriterShaLike(co),
		kennung_fmt.MakeFDCliFormat(
			co,
			u.standort.MakeRelativePathStringFormatWriter(),
		),
		u.StringFormatWriterKennung(co),
		u.StringFormatWriterTyp(co),
		u.StringFormatWriterBezeichnung(co),
		u.StringFormatWriterEtiketten(co),
	)

	return string_format_writer.MakeDelim[*sku.CheckedOut](
		"\n",
		u.Out(),
		p,
	)
}
