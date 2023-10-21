package kennung_fmt

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type kennungCliFormat struct {
	options            erworben_cli_print_options.PrintOptions
	stringFormatWriter schnittstellen.StringFormatWriter[string]
	abbr               kennung.Abbr
}

func MakeKennungCliFormat(
	options erworben_cli_print_options.PrintOptions,
	co string_format_writer.ColorOptions,
	abbr kennung.Abbr,
) *kennungCliFormat {
	return &kennungCliFormat{
		options: options,
		stringFormatWriter: string_format_writer.MakeColor[string](
			co,
			string_format_writer.MakeString[string](),
			string_format_writer.ColorTypePointer,
		),
		abbr: abbr,
	}
}

func (f *kennungCliFormat) WriteStringFormat(
	w io.StringWriter,
	k kennung.Kennung2,
) (n int64, err error) {
	if f.options.Abbreviations.Hinweisen {
		if k, err = f.abbr.AbbreviateHinweisOnly(k); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	parts := k.Parts()

	var n1 int64

	n1, err = f.stringFormatWriter.WriteStringFormat(w, parts[0])
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	var n2 int
	n2, err = w.WriteString(parts[1])
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = f.stringFormatWriter.WriteStringFormat(w, parts[2])
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
