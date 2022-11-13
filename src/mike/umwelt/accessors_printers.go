package umwelt

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/lima/store_working_directory"
)

func (u *Umwelt) PrinterZettelTransacted() collections.WriterFunc[*zettel_transacted.Zettel] {
	return collections.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelTransacted(),
	)
}

func (u *Umwelt) PrinterZettelCheckedOut() collections.WriterFunc[*zettel_checked_out.Zettel] {
	return collections.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelCheckedOut(),
	)
}

func (u *Umwelt) PrinterFileNotRecognized() collections.WriterFunc[*store_working_directory.File] {
	return collections.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatFileNotRecognized(),
	)
}

func (u *Umwelt) PrinterFileRecognized() collections.WriterFunc[*store_working_directory.FileRecognized] {
	return collections.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatFileRecognized(),
	)
}
