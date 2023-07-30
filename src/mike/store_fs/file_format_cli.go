package store_fs

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

type FileRecognized struct {
	kennung.FD
	Zettel *zettel.Transacted
}

// (recognized) [path.ext@sha]
//
//	[kopf/schwanz@sha !typ "bez"]
func MakeCliFormatRecognized(
	cw format.FuncColorWriter,
	s standort.Standort,
	sf schnittstellen.FuncWriterFormat[schnittstellen.ShaLike],
	znf schnittstellen.FuncWriterFormat[zettel.Transacted],
) schnittstellen.FuncWriterFormat[FileRecognized] {
	return func(w io.Writer, zr FileRecognized) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAligned(format.StringRecognized),
			format.MakeWriter(znf, *zr.Zettel),
			format.MakeFormatString("\n"),
			format.MakeFormatStringRightAligned("["),
			cw(s.MakeWriterRelativePath(zr.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, zr.Sha.GetShaLike()),
			format.MakeFormatString("]"),
		)
	}
}
