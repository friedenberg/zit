package string_format_writer

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
)

func MakeFunc[T any](
	f schnittstellen.FuncStringWriterFormat[T],
) schnittstellen.StringFormatWriter[T] {
	return funk[T](f)
}

type funk[T any] schnittstellen.FuncStringWriterFormat[T]

func (f funk[T]) WriteStringFormat(
	w schnittstellen.WriterAndStringWriter,
	e T,
) (int64, error) {
	return schnittstellen.FuncStringWriterFormat[T](f)(w, e)
}
