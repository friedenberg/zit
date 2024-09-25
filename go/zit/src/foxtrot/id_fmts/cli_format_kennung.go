package id_fmts

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type objectIdCliFormat struct {
	options              print_options.General
	sfwColor, sfwNoColor interfaces.StringFormatWriter[*catgut.String]
	abbr                 ids.Abbr
}

func MakeObjectIdCliFormat(
	options print_options.General,
	co string_format_writer.ColorOptions,
	abbr ids.Abbr,
) *objectIdCliFormat {
	return &objectIdCliFormat{
		options: options,
		sfwColor: string_format_writer.MakeColor(
			co,
			catgut.StringFormatWriterString,
			string_format_writer.ColorTypeId,
		),
		sfwNoColor: catgut.StringFormatWriterString,
		abbr:       abbr,
	}
}

func (f *objectIdCliFormat) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	k *ids.ObjectId,
) (n int64, err error) {
	if f.options.Abbreviations.Hinweisen {
		k1 := ids.GetObjectIdPool().Get()
		defer ids.GetObjectIdPool().Put(k1)

		if err = k1.ResetWithIdLike(k); err != nil {
			err = errors.Wrap(err)
			return
		}

		k = k1

		if err = f.abbr.AbbreviateZettelIdOnly(k1); err != nil {
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

	if parts.Middle != '\x00' {
		sm := catgut.GetPool().Get()
		defer catgut.GetPool().Put(sm)
		sm.WriteRune(rune(parts.Middle))
		n1, err = f.sfwNoColor.WriteStringFormat(w, sm)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	n1, err = f.sfwColor.WriteStringFormat(w, parts.Right)
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
