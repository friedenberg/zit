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
	fieldFormatWriter    interfaces.StringFormatWriter[string_format_writer.Field]
}

func MakeItemDeletedStringWriterFormat(
	config interfaces.Config,
	co string_format_writer.ColorOptions,
	fieldFormatWriter interfaces.StringFormatWriter[string_format_writer.Field],
) *itemDeletedStringFormatWriter {
	return &itemDeletedStringFormatWriter{
		Config:             config,
		rightAlignedWriter: string_format_writer.MakeRightAligned(),
		idStringFormatWriter: string_format_writer.MakeColor(
			co,
			string_format_writer.MakeString[string](),
			string_format_writer.ColorTypeId,
		),
		fieldFormatWriter: fieldFormatWriter,
	}
}

func (f *itemDeletedStringFormatWriter) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	item browserItem,
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

	{
		n2, err = f.fieldFormatWriter.WriteStringFormat(
			sw,
			string_format_writer.Field{
				Key:       "id",
				Value:     item.Id.String(),
				ColorType: string_format_writer.ColorTypeId,
			},
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	prefix := "\n" + string_format_writer.StringIndentWithSpace

	if item.Title != "" {
		n2, err = f.fieldFormatWriter.WriteStringFormat(
			sw,
			string_format_writer.Field{
				Key:       "title",
				Value:     item.Title,
				ColorType: string_format_writer.ColorTypeUserData,
				Prefix:    prefix,
			},
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	{
		var u *url.URL

		if u, err = item.GetUrl(); err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.fieldFormatWriter.WriteStringFormat(
			sw,
			string_format_writer.Field{
				Key:       "url",
				Value:     u.String(),
				ColorType: string_format_writer.ColorTypeUserData,
				Prefix:    prefix,
			},
		)
		n += int64(n2)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = sw.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
