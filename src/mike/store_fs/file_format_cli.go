package store_fs

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/bravo/format"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
)

// (not recognized) [path.ext@sha]
func MakeCliFormatNotRecognized(
	s standort.Standort,
	sf format.FormatWriterFunc[sha.Sha],
) format.FormatWriterFunc[File] {
	return func(w io.Writer, fu *File) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("(not recognized) [%s@", fu.Path),
			format.MakeWriter(sf, &fu.Sha),
			format.MakeFormatString("]"),
		)
	}
}

type FileRecognized struct {
	File
	Recognized zettel_transacted.MutableSet
}

// (recognized) [path.ext@sha]
//
//	[kopf/schwanz@sha !typ "bez"]
func MakeCliFormatRecognized(
	s standort.Standort,
	sf format.FormatWriterFunc[sha.Sha],
	znf format.FormatWriterFunc[zettel_named.Zettel],
) format.FormatWriterFunc[FileRecognized] {
	return func(w io.Writer, zr *FileRecognized) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatString("(recognized) [%s@", zr.Path),
			format.MakeWriter(sf, &zr.Sha),
			format.MakeFormatString("]\n"),
			func(w io.Writer) (n int64, err error) {
				err = zr.Recognized.Each(
					func(zt *zettel_transacted.Zettel) (err error) {
						var n2 int

						if n2, err = io.WriteString(w, "             "); err != nil {
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
