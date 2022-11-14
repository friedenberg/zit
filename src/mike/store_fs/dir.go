package store_fs

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/delta/standort"
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
