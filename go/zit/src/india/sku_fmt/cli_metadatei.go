package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/bezeichnung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
)

type cliMetadatei struct {
	options       erworben_cli_print_options.PrintOptions
	contentPrefix string

	writeTyp         bool
	writeBezeichnung bool
	writeEtiketten   bool

	shaStringFormatWriter         schnittstellen.StringFormatWriter[schnittstellen.ShaLike]
	typStringFormatWriter         schnittstellen.StringFormatWriter[*kennung.Typ]
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung]
	etikettenStringFormatWriter   schnittstellen.StringFormatWriter[*kennung.Etikett]
}

func MakeCliMetadateiFormatShort(
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
	typStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Typ],
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung],
	etikettenStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Etikett],
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
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
	typStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Typ],
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung],
	etikettenStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Etikett],
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
	sw schnittstellen.WriterAndStringWriter,
	o *metadatei.Metadatei,
) (n int64, err error) {
	var n1 int
	var n2 int64

	sh := &o.Akte

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
		t := o.GetMetadatei().GetTypPtr()

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
		b := o.GetMetadatei().GetBezeichnungPtr()

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
	sw schnittstellen.WriterAndStringWriter,
	o *metadatei.Metadatei,
	didWriteBezeichnung bool,
) (n int64, err error) {
	if !f.options.PrintEtikettenAlways &&
		(!f.writeEtiketten && didWriteBezeichnung) {
		return
	}

	b := o.GetMetadatei().GetEtiketten()

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
