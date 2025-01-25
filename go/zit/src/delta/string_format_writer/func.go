package string_format_writer

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func MakeFunc[T any](
	f interfaces.FuncStringWriterFormat[T],
) interfaces.StringFormatWriter[T] {
	return funk[T](f)
}

type funk[T any] interfaces.FuncStringWriterFormat[T]

func (f funk[T]) WriteStringFormat(
	e T,
	w interfaces.WriterAndStringWriter,
) (int64, error) {
	return interfaces.FuncStringWriterFormat[T](f)(w, e)
}
