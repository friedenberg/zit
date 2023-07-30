package store_fs

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type fileNotRecognizedStringWriterFormat struct {
	rightAlignedWriter    schnittstellen.StringFormatWriter[string]
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike]
	fdStringFormatWriter  schnittstellen.StringFormatWriter[*kennung.FD]
}

func MakeFileNotRecognizedStringWriterFormat(
	fdStringFormatWriter schnittstellen.StringFormatWriter[*kennung.FD],
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
) *fileNotRecognizedStringWriterFormat {
	return &fileNotRecognizedStringWriterFormat{
		rightAlignedWriter:    format.MakeRightAlignedStringFormatWriter(),
		shaStringFormatWriter: shaStringFormatWriter,
		fdStringFormatWriter:  fdStringFormatWriter,
	}
}

func (f *fileNotRecognizedStringWriterFormat) WriteStringFormat(
	sw io.StringWriter,
	fd *kennung.FD,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	prefix := format.StringUnrecognized

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

	n2, err = f.shaStringFormatWriter.WriteStringFormat(sw, fd.Sha)
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
