package string_format_writer

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func MakeBracketWrapped[T any](
	sfw schnittstellen.StringFormatWriter[T],
) schnittstellen.StringFormatWriter[T] {
	return &bracketWrapped[T]{
		stringFormatWriter: sfw,
	}
}

type bracketWrapped[T any] struct {
	stringFormatWriter schnittstellen.StringFormatWriter[T]
}

func (f bracketWrapped[T]) WriteStringFormat(
	w io.StringWriter,
	e T,
) (n int64, err error) {
	var (
		n1 int
		n2 int64
	)

	n1, err = w.WriteString("[")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n2, err = f.stringFormatWriter.WriteStringFormat(w, e)
	n += int64(n2)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	n1, err = w.WriteString("]")
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
