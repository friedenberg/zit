package id_fmts

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type shaCliFormat struct {
	ids.Abbr
	stringFormatWriter interfaces.StringFormatWriter[string]
}

func MakeShaCliFormat(
	options print_options.General,
	co string_format_writer.ColorOptions,
	abbr ids.Abbr,
) *shaCliFormat {
	return &shaCliFormat{
		Abbr: abbr,
		stringFormatWriter: string_format_writer.MakeColor[string](
			co,
			string_format_writer.MakeString[string](),
			string_format_writer.ColorTypeHash,
		),
	}
}

func (f *shaCliFormat) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	s interfaces.Sha,
) (n int64, err error) {
	v := s.String()

	if f.Abbr.Sha.Abbreviate != nil {
		var v1 string

		sh := sha.Make(s)

		if v1, err = f.Abbr.Sha.Abbreviate(sh); err != nil {
			err = errors.Wrap(err)
			return
		}

		if v1 != "" {
			v = v1
		} else {
			ui.Todo("abbreviate sha produced empty string")
		}
	}

	return f.stringFormatWriter.WriteStringFormat(w, v)
}
