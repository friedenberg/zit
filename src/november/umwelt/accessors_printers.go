package umwelt

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kasten"
	"github.com/friedenberg/zit/src/lima/zettel_checked_out"
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

func (u *Umwelt) PrinterZettelCheckedOut() schnittstellen.FuncIter[*zettel_checked_out.Zettel] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelCheckedOut(),
	)
}

func (u *Umwelt) PrinterZettelCheckedOutFresh() schnittstellen.FuncIter[*zettel_checked_out.Zettel] {
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

func (u *Umwelt) PrinterPathDeleted() schnittstellen.FuncIter[*store_fs.Dir] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatDirDeleted(),
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
