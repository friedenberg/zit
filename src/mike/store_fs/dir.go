package store_fs

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
	"github.com/friedenberg/zit/src/echo/kennung"
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
		rightAlignedWriter:   string_format_writer.MakeRightAligned(),
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

	prefix := string_format_writer.StringDeleted

	if f.dryRun {
		prefix = string_format_writer.StringWouldDelete
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
