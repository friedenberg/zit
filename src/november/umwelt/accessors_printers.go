package umwelt

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/string_writer_format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/hotel/sku_formats"
	"github.com/friedenberg/zit/src/mike/store_fs"
)

func (u *Umwelt) FormatColorOptions() (o string_writer_format.ColorOptions) {
	o.OffEntirely = !u.outIsTty
	return
}

func (u *Umwelt) StringFormatWriterShaLike() schnittstellen.StringFormatWriter[schnittstellen.ShaLike] {
	return kennung.MakeShaCliFormat(
		u.FormatColorOptions(),
		u.StoreObjekten().GetAbbrStore().Shas().Abbreviate,
	)
}

func (u *Umwelt) StringFormatWriterKennung() schnittstellen.StringFormatWriter[kennung.KennungPtr] {
	return kennung.MakeKennungCliFormat(
		u.konfig.Options,
		u.FormatColorOptions(),
		u.MakeKennungExpanders(),
	)
}

func (u *Umwelt) StringFormatWriterTyp() schnittstellen.StringFormatWriter[*kennung.Typ] {
	return kennung.MakeTypCliFormat(u.FormatColorOptions())
}

func (u *Umwelt) StringFormatWriterBezeichnung() schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung] {
	return bezeichnung.MakeCliFormat2(u.FormatColorOptions())
}

func (u *Umwelt) StringFormatWriterEtiketten() schnittstellen.StringFormatWriter[kennung.EtikettSet] {
	return kennung.MakeEtikettenCliFormat()
}

func (u *Umwelt) StringFormatWriterSkuLikePtr() schnittstellen.StringFormatWriter[sku.SkuLikePtr] {
	return sku_formats.MakeCliFormat(
		sku_formats.CliOptions{
			PrefixTai:              u.konfig.UsePrintTime(),
			AlwaysIncludeEtiketten: u.konfig.UsePrintEtiketten(),
		},
		u.StringFormatWriterShaLike(),
		u.StringFormatWriterKennung(),
		u.StringFormatWriterTyp(),
		u.StringFormatWriterBezeichnung(),
		u.StringFormatWriterEtiketten(),
	)
}

func (u *Umwelt) PrinterSkuLikePtr() schnittstellen.FuncIter[sku.SkuLikePtr] {
	sw := u.StringFormatWriterSkuLikePtr()

	return string_writer_format.MakeDelim[sku.SkuLikePtr](
		"\n",
		u.Out(),
		sw,
	)
}

func (u *Umwelt) PrinterTransactedLike() schnittstellen.FuncIter[objekte.TransactedLikePtr] {
	sw := u.StringFormatWriterSkuLikePtr()

	return string_writer_format.MakeDelim[objekte.TransactedLikePtr](
		"\n",
		u.Out(),
		string_writer_format.MakeFunc[objekte.TransactedLikePtr](
			func(w io.StringWriter, o objekte.TransactedLikePtr) (n int64, err error) {
				return sw.WriteStringFormat(w, o.GetSkuLikePtr())
			},
		),
	)
}

func (u *Umwelt) PrinterFileNotRecognized() schnittstellen.FuncIter[*kennung.FD] {
	p := store_fs.MakeFileNotRecognizedStringWriterFormat(
		kennung.MakeFDCliFormat(
			u.FormatColorOptions(),
			u.standort.MakeRelativePathStringFormatWriter(),
		),
		u.StringFormatWriterShaLike(),
	)

	return string_writer_format.MakeDelim[*kennung.FD](
		"\n",
		u.Out(),
		p,
	)
}

func (u *Umwelt) PrinterFDDeleted() schnittstellen.FuncIter[*kennung.FD] {
	p := store_fs.MakeFDDeletedStringWriterFormat(
		u.Konfig().DryRun,
		kennung.MakeFDCliFormat(
			u.FormatColorOptions(),
			u.standort.MakeRelativePathStringFormatWriter(),
		),
	)

	return string_writer_format.MakeDelim[*kennung.FD](
		"\n",
		u.Out(),
		p,
	)
}

func (u *Umwelt) PrinterHeader() schnittstellen.FuncIter[string] {
	return string_writer_format.MakeDelim[string](
		"\n",
		u.Out(),
		string_writer_format.MakeIndentedHeader(u.FormatColorOptions()),
	)
}

func (u *Umwelt) PrinterCheckedOutLike() schnittstellen.FuncIter[objekte.CheckedOutLikePtr] {
	p := objekte.MakeCliFormat(
		objekte.CliOptions{},
		u.StringFormatWriterShaLike(),
		kennung.MakeFDCliFormat(
			u.FormatColorOptions(),
			u.standort.MakeRelativePathStringFormatWriter(),
		),
		u.StringFormatWriterKennung(),
		u.StringFormatWriterTyp(),
		u.StringFormatWriterBezeichnung(),
		u.StringFormatWriterEtiketten(),
	)

	return string_writer_format.MakeDelim[objekte.CheckedOutLikePtr](
		"\n",
		u.Out(),
		p,
	)
}
