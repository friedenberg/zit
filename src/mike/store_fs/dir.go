package store_fs

import (
	"io"

	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/fd"
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
	fdw format.FormatWriterFunc[fd.FD],
) format.FormatWriterFunc[fd.FD] {
	return func(w io.Writer, fd *fd.FD) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(format.StringDeleted),
			format.MakeWriter(fdw, fd),
		)
	}
}
