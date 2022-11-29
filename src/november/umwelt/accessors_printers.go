package umwelt

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/echo/konfig"
	"github.com/friedenberg/zit/src/golf/typ"
	"github.com/friedenberg/zit/src/hotel/cwd_files"
	"github.com/friedenberg/zit/src/india/zettel"
	"github.com/friedenberg/zit/src/india/zettel_external"
	"github.com/friedenberg/zit/src/juliett/zettel_checked_out"
	"github.com/friedenberg/zit/src/mike/store_fs"
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

func (u *Umwelt) PrinterTypCheckedOut(
	verb string,
) collections.WriterFunc[*typ.External] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatTypCheckedOut(),
	)
}

func (u *Umwelt) PrinterZettelTransacted(
	verb string,
) collections.WriterFunc[*zettel.Transacted] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelTransacted(verb),
	)
}

func (u *Umwelt) PrinterZettelCheckedOut(
	mode zettel_checked_out.Mode,
) collections.WriterFunc[*zettel_checked_out.Zettel] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelCheckedOut(mode),
	)
}

func (u *Umwelt) PrinterZettelCheckedOutFresh(
	mode zettel_checked_out.Mode,
) collections.WriterFunc[*zettel_checked_out.Zettel] {
	return format.MakeWriterToWithNewLines(
		u.Out(),
		u.FormatZettelCheckedOutFresh(mode),
	)
}

func (u *Umwelt) PrinterFileNotRecognized() collections.WriterFunc[*cwd_files.File] {
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

func (u *Umwelt) PrinterFDDeleted() collections.WriterFunc[*zettel_external.FD] {
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
