package string_format_writer

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type streeng[T ~string] struct{}

func MakeString[T ~string]() schnittstellen.StringFormatWriter[T] {
	return &streeng[T]{}
}

func (f *streeng[T]) WriteStringFormat(
	sw io.StringWriter,
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
