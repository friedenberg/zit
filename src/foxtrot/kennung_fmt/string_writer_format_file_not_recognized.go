package kennung_fmt

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
	"github.com/friedenberg/zit/src/echo/fd"
)

type fileNotRecognizedStringWriterFormat struct {
	rightAlignedWriter    schnittstellen.StringFormatWriter[string]
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike]
	fdStringFormatWriter  schnittstellen.StringFormatWriter[*fd.FD]
}

func MakeFileNotRecognizedStringWriterFormat(
	fdStringFormatWriter schnittstellen.StringFormatWriter[*fd.FD],
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
) *fileNotRecognizedStringWriterFormat {
	return &fileNotRecognizedStringWriterFormat{
		rightAlignedWriter:    string_format_writer.MakeRightAligned(),
		shaStringFormatWriter: shaStringFormatWriter,
		fdStringFormatWriter:  fdStringFormatWriter,
	}
}

func (f *fileNotRecognizedStringWriterFormat) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	fd *fd.FD,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	prefix := string_format_writer.StringUnrecognized

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

	n2, err = f.shaStringFormatWriter.WriteStringFormat(sw, fd.GetShaLike())
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
