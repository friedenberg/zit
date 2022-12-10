package zettel_external

import (
	"io"

	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/golf/fd"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

// [path@sha !typ "bez"]
func MakeCliFormatZettel(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf format.FormatWriterFunc[sha.Sha],
	zf format.FormatWriterFunc[zettel.Objekte],
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(z.ZettelFD.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &z.Sku.Sha),
			format.MakeFormatString(" "),
			format.MakeWriter(zf, &z.Objekte),
			format.MakeFormatString("]"),
		)
	}
}

// [path.akte_ext@sha]
func MakeCliFormatAkte(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf format.FormatWriterFunc[sha.Sha],
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(z.AkteFD.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &z.Objekte.Akte),
			format.MakeFormatString("]"),
		)
	}
}

// [path.akte_ext@sha]
func MakeCliFormatFD(
	s standort.Standort,
	cw format.FuncColorWriter,
) format.FormatWriterFunc[fd.FD] {
	return func(w io.Writer, fd *fd.FD) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(fd.Path), format.ColorTypePointer),
			// format.MakeFormatString("@"),
			// format.MakeWriter(sf, fd.Sha),
			format.MakeFormatString("]"),
		)
	}
}
