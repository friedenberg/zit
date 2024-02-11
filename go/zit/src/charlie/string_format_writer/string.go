package string_format_writer

import (
	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
)

type streeng[T ~string] struct{}

func MakeString[T ~string]() schnittstellen.StringFormatWriter[T] {
	return &streeng[T]{}
}

func (f *streeng[T]) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
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
