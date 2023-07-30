package store_fs

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/string_writer_format"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type fdDeletedStringWriterFormat struct {
	dryRun               bool
	rightAlignedWriter   schnittstellen.StringFormatWriter[string]
	fdStringFormatWriter schnittstellen.StringFormatWriter[*kennung.FD]
}

func MakeFDDeletedStringWriterFormat(
	dryRun bool,
	fdStringFormatWriter schnittstellen.StringFormatWriter[*kennung.FD],
) *fdDeletedStringWriterFormat {
	return &fdDeletedStringWriterFormat{
		dryRun:               dryRun,
		rightAlignedWriter:   string_writer_format.MakeRightAligned(),
		fdStringFormatWriter: fdStringFormatWriter,
	}
}

func (f *fdDeletedStringWriterFormat) WriteStringFormat(
	sw io.StringWriter,
	fd *kennung.FD,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	prefix := format.StringDeleted

	if f.dryRun {
		prefix = format.StringWouldDelete
	}

	n2, err = f.rightAlignedWriter.WriteStringFormat(sw, prefix)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.fdStringFormatWriter.WriteStringFormat(
		sw,
		fd,
	)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = sw.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
