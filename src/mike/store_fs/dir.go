package store_fs

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/standort"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

// (deleted) [fd/]
func MakeCliFormatFDDeleted(
	dryRun bool,
	cw format.FuncColorWriter,
	s standort.Standort,
	fdw schnittstellen.FuncWriterFormat[kennung.FD],
) schnittstellen.FuncWriterFormat[kennung.FD] {
	prefix := format.StringDeleted

	if dryRun {
		prefix = format.StringWouldDelete
	}

	return func(w io.Writer, fd kennung.FD) (n int64, err error) {
		return format.Write(
			w,
			format.MakeFormatStringRightAligned(prefix),
			format.MakeWriter(fdw, fd),
		)
	}
}
