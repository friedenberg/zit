package id_fmts

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
)

type fileNotRecognizedStringWriterFormat struct {
	rightAlignedWriter    interfaces.StringFormatWriter[string]
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha]
	fdStringFormatWriter  interfaces.StringFormatWriter[*fd.FD]
}

func MakeFileNotRecognizedStringWriterFormat(
	fdStringFormatWriter interfaces.StringFormatWriter[*fd.FD],
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
) *fileNotRecognizedStringWriterFormat {
	return &fileNotRecognizedStringWriterFormat{
		rightAlignedWriter:    string_format_writer.MakeRightAligned(),
		shaStringFormatWriter: shaStringFormatWriter,
		fdStringFormatWriter:  fdStringFormatWriter,
	}
}

func (f *fileNotRecognizedStringWriterFormat) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
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

	n1, err = sw.WriteString("@")
	n += int64(n1)

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
