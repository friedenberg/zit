package string_format_writer

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func MakeBracketWrapped[T any](
	sfw interfaces.StringFormatWriter[T],
) interfaces.StringFormatWriter[T] {
	return &bracketWrapped[T]{
		stringFormatWriter: sfw,
	}
}

type bracketWrapped[T any] struct {
	stringFormatWriter interfaces.StringFormatWriter[T]
}

func (f bracketWrapped[T]) WriteStringFormat(
	w interfaces.WriterAndStringWriter,
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
