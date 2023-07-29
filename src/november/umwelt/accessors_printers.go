package umwelt

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/hotel/sku_formats"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/mike/store_fs"
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

func (u *Umwelt) PrinterKonfigTransacted() schnittstellen.FuncIter[*erworben.Transacted] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			erworben.MakeCliFormatTransacted(
				u.FormatColorWriter(),
				u.FormatSha(u.StoreObjekten().GetAbbrStore().Shas().Abbreviate),
			),
		),
	)
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

func (u *Umwelt) PrinterTypTransacted() schnittstellen.FuncIter[*typ.Transacted] {
	sw := sku_formats.MakeCliFormat(
		sku_formats.CliOptions{},
		u.StringFormatWriterShaLike(),
		u.StringFormatWriterKennung(),
		u.StringFormatWriterTyp(),
		u.StringFormatWriterBezeichnung(),
		u.StringFormatWriterEtiketten(),
	)

	return format.MakeDelimFuncStringFormatWriter[*typ.Transacted](
		"\n",
		u.Out(),
		format.MakeFuncStringFormatWriter(
			func(w io.StringWriter, o *typ.Transacted) (n int64, err error) {
				return sw.WriteStringFormat(w, o.GetSkuLikePtr())
			},
		),
	)
}

func (u *Umwelt) PrinterEtikettTransacted() schnittstellen.FuncIter[*etikett.Transacted] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatEtikettTransacted(),
		),
	)
}

func (u *Umwelt) PrinterKastenTransacted() schnittstellen.FuncIter[*kasten.Transacted] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatKastenTransacted(),
		),
	)
}

func (u *Umwelt) PrinterTypCheckedOut() schnittstellen.FuncIter[*typ.CheckedOut] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithCheckedOutState(
			u.FormatTypCheckedOut(),
		),
	)
}

func (u *Umwelt) PrinterKastenCheckedOut() schnittstellen.FuncIter[*kasten.CheckedOut] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithCheckedOutState(
			u.FormatKastenCheckedOut(),
		),
	)
}

func (u *Umwelt) PrinterEtikettCheckedOut() schnittstellen.FuncIter[*etikett.CheckedOut] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithCheckedOutState(
			u.FormatEtikettCheckedOut(),
		),
	)
}

func (u *Umwelt) ZettelTransactedLogPrinters() zettel.LogWriter {
	return zettel.LogWriter{
		New:       u.PrinterZettelTransactedDelta(),
		Updated:   u.PrinterZettelTransactedDelta(),
		Unchanged: u.PrinterZettelTransactedDelta(),
		Archived:  u.PrinterZettelTransactedDelta(),
	}
}

func (u *Umwelt) PrinterZettelTransacted() schnittstellen.FuncIter[*zettel.Transacted] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatZettelTransacted(),
		),
	)
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

func (u *Umwelt) PrinterZettelTransactedDelta() schnittstellen.FuncIter[*zettel.Transacted] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatZettelTransactedDelta(),
		),
	)
}

func (u *Umwelt) PrinterZettelExternal() schnittstellen.FuncIter[*zettel.External] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		u.FormatZettelExternal(),
	)
}

func (u *Umwelt) PrinterZettelCheckedOut() schnittstellen.FuncIter[*zettel.CheckedOut] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithCheckedOutState(
			u.FormatZettelCheckedOut(),
		),
	)
}

func (u *Umwelt) PrinterZettelCheckedOutFresh() schnittstellen.FuncIter[*zettel.CheckedOut] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelCheckedOut(),
	)
}

func (u *Umwelt) PrinterFileNotRecognized() schnittstellen.FuncIter[*kennung.FD] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		u.FormatFileNotRecognized(),
	)
}

func (u *Umwelt) PrinterFileRecognized() schnittstellen.FuncIter[*store_fs.FileRecognized] {
	return format.MakeWriterToWithNewLinesPtr(
		u.Out(),
		u.FormatFileRecognized(),
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

func (u *Umwelt) PrinterJustCheckedOutLike() schnittstellen.FuncIter[objekte.CheckedOutLike] {
	pzco := u.PrinterZettelExternal()
	ptco := u.PrinterTypCheckedOut()
	peco := u.PrinterEtikettCheckedOut()

	return func(co objekte.CheckedOutLike) (err error) {
		sk2 := co.GetInternalLike().GetSkuLike()

		switch sk2.GetGattung() {
		case gattung.Zettel:
			coz := co.(*zettel.CheckedOut)
			return pzco(&coz.External)

		case gattung.Typ:
			coz := co.(*typ.CheckedOut)
			return ptco(coz)

		case gattung.Etikett:
			coz := co.(*etikett.CheckedOut)
			return peco(coz)

		default:
			todo.Implement()
			_, err = fmt.Fprintf(
				u.Out(),
				"(checked out) [%s.%s]\n",
				sk2.GetKennungLike(),
				sk2.GetGattung(),
			)
		}

		return
	}
}

func (u *Umwelt) PrinterCheckedOutLike() schnittstellen.FuncIter[objekte.CheckedOutLike] {
	pzco := u.PrinterZettelCheckedOut()
	ptco := u.PrinterTypCheckedOut()
	peco := u.PrinterEtikettCheckedOut()
	pkco := u.PrinterKastenCheckedOut()

	return func(co objekte.CheckedOutLike) (err error) {
		sk2 := co.GetInternalLike().GetSkuLike()

		switch sk2.GetGattung() {
		case gattung.Zettel:
			coz := co.(*zettel.CheckedOut)
			return pzco(coz)

		case gattung.Typ:
			coz := co.(*typ.CheckedOut)
			return ptco(coz)

		case gattung.Etikett:
			coz := co.(*etikett.CheckedOut)
			return peco(coz)

		case gattung.Kasten:
			coz := co.(*kasten.CheckedOut)
			return pkco(coz)

		default:
			_, err = fmt.Fprintf(
				u.Out(),
				"(checked out) [%s.%s]\n",
				sk2.GetKennungLike(),
				sk2.GetGattung(),
			)
		}

		return
	}
}
