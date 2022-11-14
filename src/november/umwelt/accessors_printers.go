package umwelt

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	store_fs "github.com/friedenberg/zit/src/mike/store_fs"
)

func (u *Umwelt) PrinterZettelTransacted(verb string) collections.WriterFunc[*zettel_transacted.Zettel] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelTransacted(verb),
	)
}

func (u *Umwelt) PrinterZettelCheckedOut() collections.WriterFunc[*zettel_checked_out.Zettel] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelCheckedOut(),
	)
}

func (u *Umwelt) PrinterFileNotRecognized() collections.WriterFunc[*store_fs.File] {
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
