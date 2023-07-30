package umwelt

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/hotel/sku_formats"
)

func wrapWithCheckedOutState[T objekte.CheckedOutLike](
	f schnittstellen.FuncWriterFormat[T],
) schnittstellen.FuncWriterFormat[T] {
	return func(w io.Writer, e T) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAligned(e.GetState().String()),
			format.MakeWriter(f, e),
		)
	}
}

func wrapWithTimePrefixerIfNecessary[T sku.Getter](
	k schnittstellen.Konfig,
	f schnittstellen.FuncWriterFormat[T],
) schnittstellen.FuncWriterFormat[T] {
	if k.UsePrintTime() {
		return sku.MakeTimePrefixWriter(f)
	} else {
		return f
	}
}

func (u *Umwelt) StringFormatWriterShaLike() schnittstellen.StringFormatWriter[schnittstellen.ShaLike] {
	return kennung.MakeShaCliFormat2(
		u.FormatColorOptions(),
		u.StoreObjekten().GetAbbrStore().Shas().Abbreviate,
	)
}

func (u *Umwelt) StringFormatWriterKennung() schnittstellen.StringFormatWriter[kennung.KennungPtr] {
	return kennung.MakeKennungCliFormat(u.FormatColorOptions())
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

func (u *Umwelt) PrinterTransactedLike() schnittstellen.FuncIter[objekte.TransactedLikePtr] {
	sw := sku_formats.MakeCliFormat(
		sku_formats.CliOptions{PrefixTai: u.konfig.UsePrintTime()},
		u.StringFormatWriterShaLike(),
		u.StringFormatWriterKennung(),
		u.StringFormatWriterTyp(),
		u.StringFormatWriterBezeichnung(),
		u.StringFormatWriterEtiketten(),
	)

	return format.MakeDelimFuncStringFormatWriter[objekte.TransactedLikePtr](
		"\n",
		u.Out(),
		format.MakeFuncStringFormatWriter[objekte.TransactedLikePtr](
			func(w io.StringWriter, o objekte.TransactedLikePtr) (n int64, err error) {
				return sw.WriteStringFormat(w, o.GetSkuLikePtr())
			},
		),
	)
}

func (u *Umwelt) PrinterFileNotRecognized() schnittstellen.FuncIter[*kennung.FD] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		u.FormatFileNotRecognized(),
	)
}

func (u *Umwelt) PrinterFDDeleted() schnittstellen.FuncIter[*kennung.FD] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		u.FormatFDDeleted(),
	)
}

func (u *Umwelt) PrinterHeader() schnittstellen.FuncIter[*string] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		format.MakeWriterFormatStringIndentedHeader(
			u.FormatColorWriter(),
			format.StringIndent,
		),
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

	return format.MakeDelimFuncStringFormatWriter[objekte.CheckedOutLikePtr](
		"\n",
		u.Out(),
		p,
	)
}
