package string_format_writer

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
)

type quoted_streeng[T ~string] struct{}

func MakeQuotedString[T ~string]() schnittstellen.StringFormatWriter[T] {
	return &quoted_streeng[T]{}
}

func (f *quoted_streeng[T]) WriteStringFormat(
	sw schnittstellen.WriterAndStringWriter,
	e T,
) (n int64, err error) {
	var n1 int

	n1, err = fmt.Fprintf(sw, "%q", string(e))
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
