package sku_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type cli struct {
	options       erworben_cli_print_options.PrintOptions
	contentPrefix string

	kennungStringFormatWriter   interfaces.StringFormatWriter[*kennung.Kennung2]
	metadateiStringFormatWriter interfaces.StringFormatWriter[*metadatei.Metadatei]
}

func MakeCliFormatShort(
	kennungStringFormatWriter interfaces.StringFormatWriter[*kennung.Kennung2],
	metadateiStringFormatWriter interfaces.StringFormatWriter[*metadatei.Metadatei],
) *cli {
	return &cli{
		kennungStringFormatWriter:   kennungStringFormatWriter,
		metadateiStringFormatWriter: metadateiStringFormatWriter,
	}
}

func MakeCliFormat(
	options erworben_cli_print_options.PrintOptions,
	kennungStringFormatWriter interfaces.StringFormatWriter[*kennung.Kennung2],
	metadateiStringFormatWriter interfaces.StringFormatWriter[*metadatei.Metadatei],
) *cli {
	return &cli{
		options: options,
		contentPrefix: string_format_writer.StringPrefixFromOptions(
			options,
		),
		kennungStringFormatWriter:   kennungStringFormatWriter,
		metadateiStringFormatWriter: metadateiStringFormatWriter,
	}
}

func (f *cli) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	o *sku.Transacted,
) (n int64, err error) {
	var n1 int

	{
		var bracketPrefix string

		if f.options.PrintTime {
			t := o.GetTai()

			if t.IsZero() {
				ui.Err().Printf("empty tai: %s", o.GetKennung())
			} else {
				bracketPrefix = t.Format(string_format_writer.StringFormatDateTime)
			}
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

	n2, err = f.metadateiStringFormatWriter.WriteStringFormat(sw, o.GetMetadatei())
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
