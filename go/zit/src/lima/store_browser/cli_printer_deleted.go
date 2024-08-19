package store_browser

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
)

type itemDeletedStringFormatWriter struct {
	interfaces.Config
	rightAlignedWriter   interfaces.StringFormatWriter[string]
	idStringFormatWriter interfaces.StringFormatWriter[string]
}

func MakeItemDeletedStringWriterFormat(
	config interfaces.Config,
	co string_format_writer.ColorOptions,
) *itemDeletedStringFormatWriter {
	return &itemDeletedStringFormatWriter{
		Config:             config,
		rightAlignedWriter: string_format_writer.MakeRightAligned(),
		idStringFormatWriter: string_format_writer.MakeColor(
			co,
			string_format_writer.MakeString[string](),
			string_format_writer.ColorTypeId,
		),
	}
}

func (f *itemDeletedStringFormatWriter) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	co *CheckedOut,
) (n int64, err error) {
	item := co.External.item

	var (
		n1 int
		n2 int64
	)

	prefix := string_format_writer.StringDeleted

	if f.IsDryRun() {
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

	var u *url.URL

	if u, err = item.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.idStringFormatWriter.WriteStringFormat(
		sw,
		u.String(),
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
