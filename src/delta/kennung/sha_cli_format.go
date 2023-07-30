package kennung

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/bravo/string_writer_format"
)

type shaCliFormat struct {
	abbr               func(sha.Sha) (string, error)
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeShaCliFormat(
	co string_writer_format.ColorOptions,
	abbr func(sha.Sha) (string, error),
) *shaCliFormat {
	return &shaCliFormat{
		abbr: abbr,
		stringFormatWriter: string_writer_format.MakeColor[string](
			co,
			string_writer_format.MakeString[string](),
			string_writer_format.ColorTypeConstant,
		),
	}
}

func (f *shaCliFormat) WriteStringFormat(
	w io.StringWriter,
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
