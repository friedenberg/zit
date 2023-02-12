package umwelt

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/typ"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/kasten"
	"github.com/friedenberg/zit/src/kilo/zettel_external"
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

func (u *Umwelt) PrinterKonfigTransacted() collections.WriterFunc[*erworben.Transacted] {
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

func (u *Umwelt) PrinterTypTransacted() collections.WriterFunc[*typ.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatTypTransacted(),
		),
	)
}

func (u *Umwelt) PrinterEtikettTransacted() collections.WriterFunc[*etikett.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatEtikettTransacted(),
		),
	)
}

func (u *Umwelt) PrinterKastenTransacted() collections.WriterFunc[*kasten.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatKastenTransacted(),
		),
	)
}

func (u *Umwelt) PrinterTypCheckedOut() collections.WriterFunc[*typ.External] {
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

func (u *Umwelt) PrinterZettelTransacted() collections.WriterFunc[*zettel.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatZettelTransacted(),
		),
	)
}

func (u *Umwelt) PrinterZettelTransactedDelta() collections.WriterFunc[*zettel.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		wrapWithTimePrefixerIfNecessary(
			u.Konfig(),
			u.FormatZettelTransactedDelta(),
		),
	)
}

func (u *Umwelt) PrinterZettelExternal() collections.WriterFunc[*zettel_external.Zettel] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelExternal(),
	)
}

func (u *Umwelt) PrinterZettelCheckedOut() collections.WriterFunc[*zettel_checked_out.Zettel] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelCheckedOut(),
	)
}

func (u *Umwelt) PrinterZettelCheckedOutFresh() collections.WriterFunc[*zettel_checked_out.Zettel] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelCheckedOut(),
	)
}

func (u *Umwelt) PrinterFileNotRecognized() collections.WriterFunc[*kennung.FD] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatFileNotRecognized(),
	)
}

func (u *Umwelt) PrinterFileRecognized() collections.WriterFunc[*store_fs.FileRecognized] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatFileRecognized(),
	)
}

func (u *Umwelt) PrinterPathDeleted() collections.WriterFunc[*store_fs.Dir] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatDirDeleted(),
	)
}

func (u *Umwelt) PrinterFDDeleted() collections.WriterFunc[*kennung.FD] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatFDDeleted(),
	)
}

func (u *Umwelt) PrinterHeader() collections.WriterFunc[*string] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		format.MakeWriterFormatStringIndentedHeader(
			u.FormatColorWriter(),
			format.StringHeaderIndent,
		),
	)
}
