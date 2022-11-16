package zettel_external

import (
	"io"

	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

// [path@sha !typ "bez"]
func MakeCliFormatZettel(
	s standort.Standort,
	cw format.FuncColorWriter,
	sf format.FormatWriterFunc[sha.Sha],
	zf format.FormatWriterFunc[zettel.Zettel],
) format.FormatWriterFunc[Zettel] {
	return func(w io.Writer, z *Zettel) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("["),
			cw(s.MakeWriterRelativePath(z.ZettelFD.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &z.Named.Stored.Sha),
			format.MakeFormatString(" "),
			format.MakeWriter(zf, &z.Named.Stored.Zettel),
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
			format.MakeWriter(sf, &z.Named.Stored.Zettel.Akte),
			format.MakeFormatString("]"),
		)
	}
}

// [path.akte_ext@sha]
func MakeCliFormatFD(
	s standort.Standort,
	cw format.FuncColorWriter,
) format.FormatWriterFunc[FD] {
	return func(w io.Writer, fd *FD) (n int64, err error) {
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
