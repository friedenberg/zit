package sku_fmt

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
	"github.com/friedenberg/zit/src/delta/thyme"
	"github.com/friedenberg/zit/src/echo/bezeichnung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type cli struct {
	options       erworben_cli_print_options.PrintOptions
	contentPrefix string

	writeTyp         bool
	writeBezeichnung bool
	writeEtiketten   bool

	shaStringFormatWriter         schnittstellen.StringFormatWriter[schnittstellen.ShaLike]
	kennungStringFormatWriter     schnittstellen.StringFormatWriter[*kennung.Kennung2]
	typStringFormatWriter         schnittstellen.StringFormatWriter[*kennung.Typ]
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung]
	etikettenStringFormatWriter   schnittstellen.StringFormatWriter[*kennung.Etikett]
}

func MakeCliFormatShort(
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
	kennungStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Kennung2],
	typStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Typ],
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung],
	etikettenStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Etikett],
) *cli {
	return &cli{
		writeTyp:                      false,
		writeBezeichnung:              false,
		writeEtiketten:                false,
		shaStringFormatWriter:         shaStringFormatWriter,
		kennungStringFormatWriter:     kennungStringFormatWriter,
		typStringFormatWriter:         typStringFormatWriter,
		bezeichnungStringFormatWriter: bezeichnungStringFormatWriter,
		etikettenStringFormatWriter:   etikettenStringFormatWriter,
	}
}

func MakeCliFormat(
	options erworben_cli_print_options.PrintOptions,
	shaStringFormatWriter schnittstellen.StringFormatWriter[schnittstellen.ShaLike],
	kennungStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Kennung2],
	typStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Typ],
	bezeichnungStringFormatWriter schnittstellen.StringFormatWriter[*bezeichnung.Bezeichnung],
	etikettenStringFormatWriter schnittstellen.StringFormatWriter[*kennung.Etikett],
) *cli {
	return &cli{
		options: options,
		contentPrefix: string_format_writer.StringPrefixFromOptions(
			options,
		),
		writeTyp:                      true,
		writeBezeichnung:              true,
		writeEtiketten:                true,
		shaStringFormatWriter:         shaStringFormatWriter,
		kennungStringFormatWriter:     kennungStringFormatWriter,
		typStringFormatWriter:         typStringFormatWriter,
		bezeichnungStringFormatWriter: bezeichnungStringFormatWriter,
		etikettenStringFormatWriter:   etikettenStringFormatWriter,
	}
}

func (f *cli) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	o *sku.Transacted,
) (n int64, err error) {
	var n1 int

	{
		var bracketPrefix string

		if f.options.PrintTime {
			bracketPrefix = o.GetTai().Format(thyme.FormatDateTime)
		}

		if bracketPrefix != "" {
			n1, err = sw.WriteString(bracketPrefix)
			n += int64(n1)

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
	}

	n1, err = sw.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	k := &o.Kennung

	var n2 int64
	n2, err = f.kennungStringFormatWriter.WriteStringFormat(sw, k)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := o.GetAkteSha()

	if !sh.IsNull() || f.options.PrintEmptyShas {
		n1, err = sw.WriteString("@")
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n2, err = f.shaStringFormatWriter.WriteStringFormat(sw, o.GetAkteSha())
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

			n1, err = sw.WriteString("\"")
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

			n1, err = sw.WriteString("\"")
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	n2, err = f.writeStringFormatEtiketten(sw, o, didWriteBezeichnung)
	n += n2

	n1, err = sw.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f *cli) writeStringFormatEtiketten(
	sw schnittstellen.WriterAndStringWriter,
	o *sku.Transacted,
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

	for _, v := range iter.SortedValues[kennung.Etikett](b) {
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