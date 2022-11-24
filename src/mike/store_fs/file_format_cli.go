package store_fs

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/hotel/cwd_files"
	"github.com/friedenberg/zit/src/hotel/zettel"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
)

// (unrecognized) [path.ext@sha]
func MakeCliFormatNotRecognized(
	cw format.FuncColorWriter,
	s standort.Standort,
	sf format.FormatWriterFunc[sha.Sha],
) format.FormatWriterFunc[cwd_files.File] {
	return func(w io.Writer, fu *cwd_files.File) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(format.StringUnrecognized),
			format.MakeFormatString("["),
			cw(format.MakeFormatString(fu.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &fu.Sha),
			format.MakeFormatString("]"),
		)
	}
}

type FileRecognized struct {
	cwd_files.File
	Recognized zettel_transacted.MutableSet
}

// (recognized) [path.ext@sha]
//
//	[kopf/schwanz@sha !typ "bez"]
func MakeCliFormatRecognized(
	cw format.FuncColorWriter,
	s standort.Standort,
	sf format.FormatWriterFunc[sha.Sha],
	znf format.FormatWriterFunc[zettel.Named],
) format.FormatWriterFunc[FileRecognized] {
	return func(w io.Writer, zr *FileRecognized) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAlignedParen(format.StringRecognized),
			format.MakeFormatString("["),
			cw(format.MakeFormatString(zr.Path), format.ColorTypePointer),
			format.MakeFormatString("@"),
			format.MakeWriter(sf, &zr.Sha),
			format.MakeFormatString("]\n"),
			func(w io.Writer) (n int64, err error) {
				err = zr.Recognized.Each(
					func(zt *zettel_transacted.Zettel) (err error) {
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

						if n1, err = znf(w, &zt.Named); err != nil {
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
