package string_format_writer

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func MakeIndentedHeader(
	o ColorOptions,
) interfaces.StringEncoderTo[string] {
	return &indentedHeader{
		stringFormatWriter: MakeColor[string](
			o,
			MakeRightAligned(),
			ColorTypeHeading,
		),
	}
}

type indentedHeader struct {
	stringFormatWriter interfaces.StringEncoderTo[string]
}

func (f indentedHeader) EncodeStringTo(
	v string,
	w interfaces.WriterAndStringWriter,
) (n int64, err error) {
	// n1 int
	var n2 int64

	n2, err = f.stringFormatWriter.EncodeStringTo(v, w)
	n += n2

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
