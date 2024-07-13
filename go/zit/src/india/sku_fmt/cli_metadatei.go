package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/descriptions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

type cliMetadatei struct {
	options       erworben_cli_print_options.PrintOptions
	contentPrefix string

	writeTyp         bool
	writeBezeichnung bool
	writeEtiketten   bool

	shaStringFormatWriter         interfaces.StringFormatWriter[interfaces.Sha]
	typStringFormatWriter         interfaces.StringFormatWriter[*ids.Type]
	bezeichnungStringFormatWriter interfaces.StringFormatWriter[*descriptions.Description]
	etikettenStringFormatWriter   interfaces.StringFormatWriter[*ids.Tag]
}

func MakeCliMetadateiFormatShort(
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	typStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	bezeichnungStringFormatWriter interfaces.StringFormatWriter[*descriptions.Description],
	etikettenStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
) *cliMetadatei {
	return &cliMetadatei{
		writeTyp:                      false,
		writeBezeichnung:              false,
		writeEtiketten:                false,
		shaStringFormatWriter:         shaStringFormatWriter,
		typStringFormatWriter:         typStringFormatWriter,
		bezeichnungStringFormatWriter: bezeichnungStringFormatWriter,
		etikettenStringFormatWriter:   etikettenStringFormatWriter,
	}
}

func MakeCliMetadateiFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter interfaces.StringFormatWriter[interfaces.Sha],
	typStringFormatWriter interfaces.StringFormatWriter[*ids.Type],
	bezeichnungStringFormatWriter interfaces.StringFormatWriter[*descriptions.Description],
	etikettenStringFormatWriter interfaces.StringFormatWriter[*ids.Tag],
) *cliMetadatei {
	return &cliMetadatei{
		options: options,
		contentPrefix: string_format_writer.StringPrefixFromOptions(
			options,
		),
		writeTyp:                      true,
		writeBezeichnung:              true,
		writeEtiketten:                true,
		shaStringFormatWriter:         shaStringFormatWriter,
		typStringFormatWriter:         typStringFormatWriter,
		bezeichnungStringFormatWriter: bezeichnungStringFormatWriter,
		etikettenStringFormatWriter:   etikettenStringFormatWriter,
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
		t := o.GetMetadata().GetTypPtr()

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

	didWriteBezeichnung := false
	if f.writeBezeichnung {
		b := o.GetMetadata().GetBezeichnungPtr()

		if !b.IsEmpty() {
			didWriteBezeichnung = true

			n1, err = sw.WriteString(f.contentPrefix)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n2, err = f.bezeichnungStringFormatWriter.WriteStringFormat(sw, b)
			n += n2

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	n2, err = f.writeStringFormatEtiketten(sw, o, didWriteBezeichnung)
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
	didWriteBezeichnung bool,
) (n int64, err error) {
	if !f.options.PrintEtikettenAlways &&
		(!f.writeEtiketten && didWriteBezeichnung) {
		return
	}

	b := o.GetMetadata().GetEtiketten()

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
