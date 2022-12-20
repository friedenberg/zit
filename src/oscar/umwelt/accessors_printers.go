package umwelt

import (
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/etikett"
	"github.com/friedenberg/zit/src/golf/fd"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/india/typ"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/zettel_checked_out"
	"github.com/friedenberg/zit/src/november/store_fs"
)

func (u *Umwelt) PrinterKonfigTransacted(
	verb string,
) collections.WriterFunc[*konfig.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		konfig.MakeCliFormatTransacted(
			u.FormatColorWriter(),
			u.FormatSha(),
			verb,
		),
	)
}

func (u *Umwelt) PrinterTypTransacted(
	verb string,
) collections.WriterFunc[*typ.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatTypTransacted(verb),
	)
}

func (u *Umwelt) PrinterEtikettTransacted(
	verb string,
) collections.WriterFunc[*etikett.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatEtikettTransacted(verb),
	)
}

func (u *Umwelt) PrinterTypCheckedOut(
	verb string,
) collections.WriterFunc[*typ.External] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatTypCheckedOut(),
	)
}

// TODO-P0 rename
func (u *Umwelt) WriterZettelTransacted() zettel.TransactedWriters {
	return zettel.TransactedWriters{
		New:       u.PrinterZettelTransactedDelta(format.StringNew),
		Updated:   u.PrinterZettelTransactedDelta(format.StringUpdated),
		Unchanged: u.PrinterZettelTransactedDelta(format.StringUnchanged),
		Archived:  u.PrinterZettelTransactedDelta(format.StringArchived),
	}
}

func (u *Umwelt) PrinterZettelTransacted() collections.WriterFunc[*zettel.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelTransacted(),
	)
}

func (u *Umwelt) PrinterZettelTransactedDelta(
	verb string,
) collections.WriterFunc[*zettel.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelTransactedDelta(verb),
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

func (u *Umwelt) PrinterFileNotRecognized() collections.WriterFunc[*fd.FD] {
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

func (u *Umwelt) PrinterFDDeleted() collections.WriterFunc[*fd.FD] {
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
