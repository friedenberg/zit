package kennung_fmt

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/src/echo/kennung"
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

	n1, err = f.sfwColor.WriteStringFormat(w, parts.Left)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	sm := catgut.GetPool().Get()
	defer catgut.GetPool().Put(sm)
	sm.WriteRune(rune(parts.Middle))
	n1, err = f.sfwNoColor.WriteStringFormat(w, sm)
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = f.sfwColor.WriteStringFormat(w, parts.Right)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
