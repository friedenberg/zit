package id_fmts

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
)

type fdCliFormat struct {
	stringFormatWriter interfaces.StringFormatWriter[string]
}

func MakeFDCliFormat(
	co string_format_writer.ColorOptions,
	relativePathStringFormatWriter interfaces.StringFormatWriter[string],
) *fdCliFormat {
	return &fdCliFormat{
		stringFormatWriter: string_format_writer.MakeColor[string](
			co,
			relativePathStringFormatWriter,
			string_format_writer.ColorTypeId,
		),
	}
}

func (f *fdCliFormat) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
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
