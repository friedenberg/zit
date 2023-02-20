package umwelt

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/kasten"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/mike/store_fs"
)

func wrapWithTimePrefixerIfNecessary[T sku.DataIdentityGetter](
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
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			erworben.MakeCliFormatTransacted(
				u.FormatColorWriter(),
				u.FormatSha(u.StoreObjekten().GetAbbrStore().AbbreviateSha),
			),
		),
	)
}

func (u *Umwelt) PrinterTypTransacted() schnittstellen.FuncIter[*typ.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatTypTransacted(),
		),
	)
}

func (u *Umwelt) PrinterEtikettTransacted() schnittstellen.FuncIter[*etikett.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatEtikettTransacted(),
		),
	)
}

func (u *Umwelt) PrinterKastenTransacted() schnittstellen.FuncIter[*kasten.Transacted] {
	return format.MakeWriterToWithNewLines(
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
		u.FormatTypCheckedOut(),
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
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatZettelTransacted(),
		),
	)
}

func (u *Umwelt) PrinterTransactedLike() schnittstellen.FuncIter[objekte.TransactedLike] {
	z := u.FormatZettelTransacted()

	return format.MakeWriterToWithNewLines2(
		u.Out(),
		wrapWithTimePrefixerIfNecessary[objekte.TransactedLike](
			u.Konfig(),
			func(out io.Writer, tl objekte.TransactedLike) (n int64, err error) {
				switch atl := tl.(type) {
				case zettel.Transacted:
					return z(out, atl)

				case *zettel.Transacted:
					return z(out, *atl)

				default:
					err = errors.Implement()
					return
				}
			},
		),
	)
}

func (u *Umwelt) PrinterZettelTransactedDelta() schnittstellen.FuncIter[*zettel.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatZettelTransactedDelta(),
		),
	)
}

func (u *Umwelt) PrinterZettelExternal() schnittstellen.FuncIter[*zettel.External] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelExternal(),
	)
}

func (u *Umwelt) PrinterZettelCheckedOut() schnittstellen.FuncIter[*zettel.CheckedOut] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelCheckedOut(),
	)
}

func (u *Umwelt) PrinterZettelCheckedOutFresh() schnittstellen.FuncIter[*zettel.CheckedOut] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelCheckedOut(),
	)
}

func (u *Umwelt) PrinterFileNotRecognized() schnittstellen.FuncIter[*kennung.FD] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatFileNotRecognized(),
	)
}

func (u *Umwelt) PrinterFileRecognized() schnittstellen.FuncIter[*store_fs.FileRecognized] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatFileRecognized(),
	)
}

func (u *Umwelt) PrinterFDDeleted() schnittstellen.FuncIter[*kennung.FD] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatFDDeleted(),
	)
}

func (u *Umwelt) PrinterHeader() schnittstellen.FuncIter[*string] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		format.MakeWriterFormatStringIndentedHeader(
			u.FormatColorWriter(),
			format.StringHeaderIndent,
		),
	)
}

func (u *Umwelt) PrinterJustCheckedOutLike() schnittstellen.FuncIter[objekte.CheckedOutLike] {
	pzco := u.PrinterZettelExternal()
	ptco := u.PrinterTypCheckedOut()

	return func(co objekte.CheckedOutLike) (err error) {
		sk2 := co.GetInternal().GetSku2()

		switch sk2.Gattung {
		case gattung.Zettel:
			coz := co.(*zettel.CheckedOut)
			return pzco(&coz.External)

		case gattung.Typ:
			coz := co.(typ.CheckedOut)
			return ptco(&coz)

		default:
			_, err = fmt.Fprintf(
				u.Out(),
				"(checked out) [%s.%s]\n",
				sk2.Kennung,
				sk2.Gattung,
			)
		}

		return
	}
}

func (u *Umwelt) PrinterCheckedOutLike() schnittstellen.FuncIter[objekte.CheckedOutLike] {
	pzco := u.PrinterZettelCheckedOut()
	ptco := u.PrinterTypCheckedOut()

	return func(co objekte.CheckedOutLike) (err error) {
		sk2 := co.GetInternal().GetSku2()

		switch sk2.Gattung {
		case gattung.Zettel:
			coz := co.(*zettel.CheckedOut)
			return pzco(coz)

		case gattung.Typ:
			coz := co.(*typ.CheckedOut)
			return ptco(coz)

		default:
			_, err = fmt.Fprintf(
				u.Out(),
				"(checked out) [%s.%s]\n",
				sk2.Kennung,
				sk2.Gattung,
			)
		}

		return
	}
}
