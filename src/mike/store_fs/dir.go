package store_fs

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/india/zettel_external"
)

type Dir struct {
	Path string
}

func (d Dir) String() string {
	return d.Path
}

// (deleted) [dir/]
func MakeCliFormatDirDeleted(
	cw format.FuncColorWriter,
	s standort.Standort,
) format.FormatWriterFunc[Dir] {
	return func(w io.Writer, d *Dir) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(format.StringDeleted),
			format.MakeFormatString("["),
			cw(format.MakeFormatString(d.Path), format.ColorTypePointer),
			format.MakeFormatString("]"),
		)
	}
}

// (deleted) [fd/]
func MakeCliFormatFDDeleted(
	cw format.FuncColorWriter,
	s standort.Standort,
	fdw format.FormatWriterFunc[zettel_external.FD],
) format.FormatWriterFunc[zettel_external.FD] {
	return func(w io.Writer, fd *zettel_external.FD) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(format.StringDeleted),
			format.MakeWriter(fdw, fd),
		)
	}
}
