package object_metadata_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

type cliMetadatei struct {
	options       print_options.General
	contentPrefix string

	writeTyp         bool
	writeDescription bool
	writeEtiketten   bool

	Sha    interfaces.StringFormatWriter[interfaces.Sha]
	Type   interfaces.StringFormatWriter[*ids.Type]
	Fields interfaces.StringFormatWriter[string_format_writer.Box]
	Tags   interfaces.StringFormatWriter[*ids.Tag]
}

func MakeCliMetadateiFormat(
	options print_options.General,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	typStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	fieldsFormatWriter interfaces.StringFormatWriter[string_format_writer.Box],
	etikettenStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
) *cliMetadatei {
	return &cliMetadatei{
		options: options,
		contentPrefix: string_format_writer.StringPrefixFromOptions(
			options,
		),
		writeTyp:         true,
		writeDescription: true,
		writeEtiketten:   true,
		Sha:              shaStringFormatWriter,
		Type:             typStringFormatWriter,
		Fields:           fieldsFormatWriter,
		Tags:             etikettenStringFormatWriter,
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
		n1, err = sw.WriteString(" @")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.Sha.WriteStringFormat(sw, sh)
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

			n2, err = f.Type.WriteStringFormat(sw, t)
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

			n2, err = f.Fields.WriteStringFormat(
				sw,
				string_format_writer.Box{
					Contents: []string_format_writer.Field{
						{
							Value:     b.String(),
							ColorType: string_format_writer.ColorTypeUserData,
						},
					},
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

		n2, err = f.Tags.WriteStringFormat(sw, &v)
		n += n2

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
