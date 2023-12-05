package kennung_fmt

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type kennungCliFormat struct {
	options              erworben_cli_print_options.PrintOptions
	sfwColor, sfwNoColor schnittstellen.StringFormatWriter[*catgut.String]
	abbr                 kennung.Abbr
}

func MakeKennungCliFormat(
	options erworben_cli_print_options.PrintOptions,
	co string_format_writer.ColorOptions,
	abbr kennung.Abbr,
) *kennungCliFormat {
	return &kennungCliFormat{
		options: options,
		sfwColor: string_format_writer.MakeColor[*catgut.String](
			co,
			catgut.StringFormatWriter,
			string_format_writer.ColorTypePointer,
		),
		sfwNoColor: catgut.StringFormatWriter,
		abbr:       abbr,
	}
}

func (f *kennungCliFormat) WriteStringFormat(
	w schnittstellen.WriterAndStringWriter,
	k *kennung.Kennung2,
) (n int64, err error) {
	if f.options.Abbreviations.Hinweisen {
		k1 := kennung.GetKennungPool().Get()
		defer kennung.GetKennungPool().Put(k1)

		if err = k1.ResetWithKennung(k); err != nil {
			err = errors.Wrap(err)
			return
		}

		k = k1

		if err = f.abbr.AbbreviateHinweisOnly(k1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	parts := k.PartsStrings()

	var n1 int64

	n1, err = f.sfwColor.WriteStringFormat(w, parts[0])
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = f.sfwNoColor.WriteStringFormat(w, parts[1])
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = f.sfwColor.WriteStringFormat(w, parts[2])
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
