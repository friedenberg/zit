package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

type cliMetadatei struct {
	options       erworben_cli_print_options.PrintOptions
	contentPrefix string

	writeTyp         bool
	writeDescription bool
	writeEtiketten   bool

	shaStringFormatWriter       interfaces.StringFormatWriter[interfaces.Sha]
	typStringFormatWriter       interfaces.StringFormatWriter[*ids.Type]
	fieldStringFormatWriter     interfaces.StringFormatWriter[string_format_writer.Field]
	etikettenStringFormatWriter interfaces.StringFormatWriter[*ids.Tag]
}

func MakeCliMetadateiFormatShort(
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	typStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	fieldFormatWriter interfaces.StringFormatWriter[string_format_writer.Field],
	etikettenStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
) *cliMetadatei {
	return &cliMetadatei{
		writeTyp:                    false,
		writeDescription:            false,
		writeEtiketten:              false,
		shaStringFormatWriter:       shaStringFormatWriter,
		typStringFormatWriter:       typStringFormatWriter,
		fieldStringFormatWriter:     fieldFormatWriter,
		etikettenStringFormatWriter: etikettenStringFormatWriter,
	}
}

func MakeCliMetadateiFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	typStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	fieldFormatWriter interfaces.StringFormatWriter[string_format_writer.Field],
	etikettenStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
) *cliMetadatei {
	return &cliMetadatei{
		options: options,
		contentPrefix: string_format_writer.StringPrefixFromOptions(
			options,
		),
		writeTyp:                    true,
		writeDescription:            true,
		writeEtiketten:              true,
		shaStringFormatWriter:       shaStringFormatWriter,
		typStringFormatWriter:       typStringFormatWriter,
		fieldStringFormatWriter:     fieldFormatWriter,
		etikettenStringFormatWriter: etikettenStringFormatWriter,
	}
}

func (f *cliMetadatei) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	o *object_metadata.Metadata,
) (n int64, err error) {
	var n1 int
	var n2 int64

	sh := &o.Blob

	if !sh.IsNull() || f.options.PrintEmptyShas {
		n1, err = sw.WriteString("@")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.shaStringFormatWriter.WriteStringFormat(sw, sh)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if f.writeTyp {
		t := o.GetMetadata().GetTypePtr()

		if len(t.String()) > 0 {
			n1, err = sw.WriteString(f.contentPrefix)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n1, err = sw.WriteString("!")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.typStringFormatWriter.WriteStringFormat(sw, t)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	didWriteDescription := false
	if f.writeDescription {
		b := &o.Description

		if !b.IsEmpty() {
			didWriteDescription = true

			n1, err = sw.WriteString(f.contentPrefix)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.fieldStringFormatWriter.WriteStringFormat(
				sw,
				string_format_writer.Field{
					Value:     b.String(),
					ColorType: string_format_writer.ColorTypeUserData,
				},
			)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	n2, err = f.writeStringFormatEtiketten(sw, o, didWriteDescription)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *cliMetadatei) writeStringFormatEtiketten(
	sw interfaces.WriterAndStringWriter,
	o *object_metadata.Metadata,
	didWriteDescription bool,
) (n int64, err error) {
	if !f.options.PrintTagsAlways &&
		(!f.writeEtiketten && didWriteDescription) {
		return
	}

	b := o.GetMetadata().GetTags()

	if b.Len() == 0 {
		return
	}

	var n1 int
	var n2 int64

	for _, v := range iter.SortedValues(b) {
		n1, err = sw.WriteString(f.contentPrefix)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.etikettenStringFormatWriter.WriteStringFormat(sw, &v)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
