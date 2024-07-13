package kennung_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type kennungCliFormat struct {
	options              erworben_cli_print_options.PrintOptions
	sfwColor, sfwNoColor interfaces.StringFormatWriter[*catgut.String]
	abbr                 kennung.Abbr
}

func MakeKennungCliFormat(
	options erworben_cli_print_options.PrintOptions,
	co string_format_writer.ColorOptions,
	abbr kennung.Abbr,
) *kennungCliFormat {
	return &kennungCliFormat{
		options: options,
		sfwColor: string_format_writer.MakeColor(
			co,
			catgut.StringFormatWriter,
			string_format_writer.ColorTypePointer,
		),
		sfwNoColor: catgut.StringFormatWriter,
		abbr:       abbr,
	}
}

func (f *kennungCliFormat) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	k *kennung.Id,
) (n int64, err error) {
	if f.options.Abbreviations.Hinweisen {
		k1 := kennung.GetIdPool().Get()
		defer kennung.GetIdPool().Put(k1)

		if err = k1.ResetWithIdLike(k); err != nil {
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
