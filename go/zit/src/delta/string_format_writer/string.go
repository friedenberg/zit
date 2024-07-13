package string_format_writer

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

type streeng[T ~string] struct{}

func MakeString[T ~string]() interfaces.StringFormatWriter[T] {
	return &streeng[T]{}
}

func (f *streeng[T]) WriteStringFormat(
	sw interfaces.WriterAndStringWriter,
	e T,
) (n int64, err error) {
	var n1 int

	n1, err = sw.WriteString(string(e))
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
