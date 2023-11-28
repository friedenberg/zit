package kennung_fmt

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/charlie/string_format_writer"
)

type shaCliFormat struct {
	abbr               func(sha.Sha) (string, error)
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeShaCliFormat(
	options erworben_cli_print_options.PrintOptions,
	co string_format_writer.ColorOptions,
	abbr func(sha.Sha) (string, error),
) *shaCliFormat {
	if !options.Abbreviations.Shas {
		abbr = nil
	}

	return &shaCliFormat{
		abbr: abbr,
		stringFormatWriter: string_format_writer.MakeColor[string](
			co,
			string_format_writer.MakeString[string](),
			string_format_writer.ColorTypeConstant,
		),
	}
}

func (f *shaCliFormat) WriteStringFormat(
	w schnittstellen.WriterAndStringWriter,
	s schnittstellen.ShaLike,
) (n int64, err error) {
	v := s.String()

	if f.abbr != nil {
		var v1 string

		if v1, err = f.abbr(sha.Make(s)); err != nil {
			err = errors.Wrap(err)
			return
		}

		if v1 != "" {
			v = v1
		} else {
			errors.Todo("abbreviate sha produced empty string")
		}
	}

	return f.stringFormatWriter.WriteStringFormat(w, v)
}
