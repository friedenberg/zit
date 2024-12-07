package string_format_writer

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func MakeIndentedHeader(
	o ColorOptions,
) interfaces.StringFormatWriter[string] {
	return &indentedHeader{
		stringFormatWriter: MakeColor[string](
			o,
			MakeRightAligned(),
			ColorTypeHeading,
		),
	}
}

type indentedHeader struct {
	stringFormatWriter interfaces.StringFormatWriter[string]
}

func (f indentedHeader) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
	v string,
) (n int64, err error) {
	// n1 int
	var n2 int64

	n2, err = f.stringFormatWriter.WriteStringFormat(w, v)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
