package store_browser

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
)

type itemDeletedStringFormatWriter struct {
	interfaces.Config
	rightAlignedWriter   interfaces.StringFormatWriter[string]
	idStringFormatWriter interfaces.StringFormatWriter[string]
	fieldsFormatWriter   interfaces.StringFormatWriter[string_format_writer.Fields]
}

func MakeItemDeletedStringWriterFormat(
	config interfaces.Config,
	co string_format_writer.ColorOptions,
	fieldsFormatWriter interfaces.StringFormatWriter[string_format_writer.Fields],
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
	item Item,
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

	fields := []string_format_writer.Field{
		{
			Key:       "id",
			Value:     item.Id.String(),
			ColorType: string_format_writer.ColorTypeId,
		},
	}

	prefix := "\n" + string_format_writer.StringIndentWithSpace

	if item.Title != "" {
		fields = append(
			fields,
			string_format_writer.Field{
				Key:       "title",
				Value:     item.Title,
				ColorType: string_format_writer.ColorTypeUserData,
				Prefix:    prefix,
			},
		)
	}

	fields = append(
		fields,
		string_format_writer.Field{
			Key:       "url",
			Value:     item.Url.String(),
			ColorType: string_format_writer.ColorTypeUserData,
			Prefix:    prefix,
		},
	)

	n2, err = f.fieldsFormatWriter.WriteStringFormat(
		sw,
		string_format_writer.Fields{Boxed: fields},
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
