package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type itemDeletedStringFormatWriter struct {
	interfaces.Config
	rightAlignedWriter   interfaces.StringFormatWriter[string]
	idStringFormatWriter interfaces.StringFormatWriter[string]
	fieldsFormatWriter   interfaces.StringFormatWriter[string_format_writer.Box]
}

func MakeItemDeletedStringWriterFormat(
	config interfaces.Config,
	co string_format_writer.ColorOptions,
	fieldsFormatWriter interfaces.StringFormatWriter[string_format_writer.Box],
) *itemDeletedStringFormatWriter {
	return &itemDeletedStringFormatWriter{
		Config:             config,
		rightAlignedWriter: string_format_writer.MakeRightAligned(),
		idStringFormatWriter: string_format_writer.MakeColor(
			co,
			string_format_writer.MakeString[string](),
			string_format_writer.ColorTypeId,
		),
		fieldsFormatWriter: fieldsFormatWriter,
	}
}

func (f *itemDeletedStringFormatWriter) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	o *sku.Transacted,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	prefixOne := string_format_writer.StringDeleted

	if f.IsDryRun() {
		prefixOne = string_format_writer.StringWouldDelete
	}

	n2, err = f.rightAlignedWriter.WriteStringFormat(sw, prefixOne)
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

	n2, err = f.fieldsFormatWriter.WriteStringFormat(
		sw,
		string_format_writer.Box{Contents: o.Metadata.Fields},
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
