package store_fs

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/juliett/zettel"
)

// (unrecognized) [path.ext@sha]
func MakeCliFormatNotRecognized(
	cw format.FuncColorWriter,
	s standort.Standort,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
) schnittstellen.FuncWriterFormat[kennung.FD] {
	return func(w io.Writer, fu kennung.FD) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(format.StringUnrecognized),
			format.MakeFormatString("["),
			cw(format.MakeFormatString(fu.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, fu.Sha.GetSha()),
			format.MakeFormatString("]"),
		)
	}
}

type FileRecognized struct {
	kennung.FD
	Recognized zettel.MutableSet
}

// (recognized) [path.ext@sha]
//
//	[kopf/schwanz@sha !typ "bez"]
func MakeCliFormatRecognized(
	cw format.FuncColorWriter,
	s standort.Standort,
	sf schnittstellen.FuncWriterFormat[schnittstellen.Sha],
	znf schnittstellen.FuncWriterFormat[zettel.Objekte],
) schnittstellen.FuncWriterFormat[FileRecognized] {
	return func(w io.Writer, zr FileRecognized) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(format.StringRecognized),
			format.MakeFormatString("["),
			cw(format.MakeFormatString(zr.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, zr.Sha.GetSha()),
			format.MakeFormatString("]\n"),
			func(w io.Writer) (n int64, err error) {
				err = zr.Recognized.Each(
					func(zt *zettel.Transacted) (err error) {
						var n2 int

						if n2, err = io.WriteString(
							w,
							strings.Repeat(" ", format.LenStringMax),
						); err != nil {
							err = errors.Wrap(err)
							return
						}

						n += int64(n2)

						var n1 int64

						if n1, err = znf(w, zt.Objekte); err != nil {
							err = errors.Wrap(err)
							return
						}

						n += n1

						return
					},
				)

				return
			},
		)
	}
}
