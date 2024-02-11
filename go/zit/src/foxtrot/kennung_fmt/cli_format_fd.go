package kennung_fmt

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/string_format_writer"
	"code.linenisgreat.com/zit/src/echo/fd"
)

type fdCliFormat struct {
	stringFormatWriter schnittstellen.StringFormatWriter[string]
}

func MakeFDCliFormat(
	co string_format_writer.ColorOptions,
	relativePathStringFormatWriter schnittstellen.StringFormatWriter[string],
) *fdCliFormat {
	return &fdCliFormat{
		stringFormatWriter: string_format_writer.MakeColor[string](
			co,
			relativePathStringFormatWriter,
			string_format_writer.ColorTypePointer,
		),
	}
}

func (f *fdCliFormat) WriteStringFormat(
	w schnittstellen.WriterAndStringWriter,
	k *fd.FD,
) (n int64, err error) {
	// TODO-P2 add abbreviation

	var n1 int64

	n1, err = f.stringFormatWriter.WriteStringFormat(w, k.String())
	n += n1

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
