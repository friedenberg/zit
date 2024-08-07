package browser

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type cliCheckedOut struct {
	options erworben_cli_print_options.PrintOptions

	rightAlignedWriter         interfaces.StringFormatWriter[string]
	shaStringFormatWriter      interfaces.StringFormatWriter[interfaces.Sha]
	objectIdStringFormatWriter interfaces.StringFormatWriter[*ids.ObjectId]
	metadataStringFormatWriter interfaces.StringFormatWriter[*object_metadata.Metadata]

	typeStringFormatWriter        interfaces.StringFormatWriter[*ids.Type]
	descriptionStringFormatWriter interfaces.StringFormatWriter[*descriptions.Description]
	tagsStringFormatWriter        interfaces.StringFormatWriter[*ids.Tag]

	fieldFormatWriter interfaces.StringFormatWriter[string_format_writer.Field]
}

func MakeCliCheckedOutFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	objectIdStringFormatWriter interfaces.StringFormatWriter[*ids.ObjectId],
	metadataStringFormatWriter interfaces.StringFormatWriter[*object_metadata.Metadata],
	typeStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	descriptionStringFormatWriter interfaces.StringFormatWriter[*descriptions.Description],
	tagsStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
	fieldFormatWriter interfaces.StringFormatWriter[string_format_writer.Field],
) *cliCheckedOut {
	return &cliCheckedOut{
		options:                       options,
		rightAlignedWriter:            string_format_writer.MakeRightAligned(),
		shaStringFormatWriter:         shaStringFormatWriter,
		objectIdStringFormatWriter:    objectIdStringFormatWriter,
		metadataStringFormatWriter:    metadataStringFormatWriter,
		typeStringFormatWriter:        typeStringFormatWriter,
		descriptionStringFormatWriter: descriptionStringFormatWriter,
		tagsStringFormatWriter:        tagsStringFormatWriter,
		fieldFormatWriter:             fieldFormatWriter,
	}
}

func (f *cliCheckedOut) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	col sku.CheckedOutLike,
) (n int64, err error) {
	co := col.(*CheckedOut)
	var (
		n1 int
		n2 int64
	)

	{
		var stateString string

		if co.State == checked_out_state.Error {
			stateString = co.Error.Error()
		} else {
			stateString = co.State.String()
		}

		n2, err = f.rightAlignedWriter.WriteStringFormat(sw, stateString)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if co.State != checked_out_state.Untracked {
		n2, err = f.objectIdStringFormatWriter.WriteStringFormat(
			sw,
			&co.Internal.ObjectId,
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.metadataStringFormatWriter.WriteStringFormat(
			sw,
			&co.Internal.Metadata,
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = sw.WriteString("\n")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.rightAlignedWriter.WriteStringFormat(
			sw,
			"",
		)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = sw.WriteString(" ")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	item := co.External.item
	browser := &co.External.browser

	prefix := "\n" + string_format_writer.StringIndentWithSpace

	{
		{
			n2, err = f.fieldFormatWriter.WriteStringFormat(
				sw,
				string_format_writer.Field{
					Key:       "id",
					Value:     co.External.ObjectId.String(),
					ColorType: string_format_writer.ColorTypeId,
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
			n2, err = f.fieldFormatWriter.WriteStringFormat(
				sw,
				string_format_writer.Field{
					Key:       "type",
					Value:     browser.Metadata.Type.String(),
					ColorType: string_format_writer.ColorTypeType,
					Prefix:    prefix,
				},
			)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		if !browser.Metadata.Description.IsEmpty() {
			n2, err = f.fieldFormatWriter.WriteStringFormat(
				sw,
				string_format_writer.Field{
					Key:       "title",
					Value:     browser.Metadata.Description.String(),
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
	}

	{
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

		tags := browser.Metadata.GetTags()
		first := true

		if tags.Len() > 0 {
			for _, v := range iter.SortedValues(browser.Metadata.GetTags()) {
				field := string_format_writer.Field{
					Value:  v.String(),
					Prefix: " ",
				}

				if first {
					field.Prefix = prefix
				}

				n2, err = f.fieldFormatWriter.WriteStringFormat(sw, field)
				n += int64(n2)

				if err != nil {
					err = errors.Wrap(err)
					return
				}

				first = false
			}
		}
	}

	n1, err = sw.WriteString("\n")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.rightAlignedWriter.WriteStringFormat(sw, "")
	n += int64(n2)

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
