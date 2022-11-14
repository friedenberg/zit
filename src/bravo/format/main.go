package format

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type WriterFunc func(io.Writer) (int64, error)
type FormatWriterFunc[T any] func(io.Writer, *T) (int64, error)
type FormatReaderFunc[T any] func(io.Reader, *T) (int64, error)
type FuncColorWriter func(WriterFunc, ColorType) WriterFunc

func MakeWriter[T any](
	wff FormatWriterFunc[T],
	e *T,
) WriterFunc {
	return func(w io.Writer) (int64, error) {
		return wff(w, e)
	}
}

func MakeFormatString(
	f string,
	vs ...interface{},
) WriterFunc {
	return func(w io.Writer) (n int64, err error) {
		var n1 int

		if n1, err = io.WriteString(w, fmt.Sprintf(f, vs...)); err != nil {
			n = int64(n1)
			err = errors.Wrap(err)
			return
		}

		n = int64(n1)

		return
	}
}

func MakeFormatStringer[T fmt.Stringer]() FormatWriterFunc[T] {
	return func(w io.Writer, e *T) (n int64, err error) {
		var n1 int

		if n1, err = io.WriteString(w, T(*e).String()); err != nil {
			n = int64(n1)
			err = errors.Wrap(err)
			return
		}

		n = int64(n1)

		return
	}
}
