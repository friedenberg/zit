package store_fs

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/foxtrot/fd"
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
) schnittstellen.FuncWriterFormat[Dir] {
	return func(w io.Writer, d Dir) (n int64, err error) {
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
	fdw schnittstellen.FuncWriterFormat[fd.FD],
) schnittstellen.FuncWriterFormat[fd.FD] {
	return func(w io.Writer, fd fd.FD) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(format.StringDeleted),
			format.MakeWriter(fdw, fd),
		)
	}
}
