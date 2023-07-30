package string_writer_format

import (
	"io"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func MakeFunc[T any](
	f schnittstellen.FuncStringWriterFormat[T],
) schnittstellen.StringFormatWriter[T] {
	return funk[T](f)
}

type funk[T any] schnittstellen.FuncStringWriterFormat[T]

func (f funk[T]) WriteStringFormat(
	w io.StringWriter,
	e T,
) (int64, error) {
	return schnittstellen.FuncStringWriterFormat[T](f)(w, e)
}
