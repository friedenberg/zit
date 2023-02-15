package format

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type FuncColorWriter func(schnittstellen.FuncWriter, ColorType) schnittstellen.FuncWriter

func MakeWriter[T any](
	wff schnittstellen.FuncWriterFormat[T],
	e T,
) schnittstellen.FuncWriter {
	return func(w io.Writer) (int64, error) {
		return wff(w, e)
	}
}

func MakeWriterPtr[T any](
	wff schnittstellen.FuncWriterFormat[*T],
	e *T,
) schnittstellen.FuncWriter {
	return func(w io.Writer) (int64, error) {
		return wff(w, e)
	}
}

func MakeFormatString(
	f string,
	vs ...interface{},
) schnittstellen.FuncWriter {
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

func MakeStringer(
	v fmt.Stringer,
) schnittstellen.FuncWriter {
	return func(w io.Writer) (n int64, err error) {
		var n1 int

		if n1, err = io.WriteString(w, v.String()); err != nil {
			n = int64(n1)
			err = errors.Wrap(err)
			return
		}

		n = int64(n1)

		return
	}
}

func MakeFormatStringer[T schnittstellen.ValueLike](
	sf schnittstellen.FuncString[schnittstellen.Set[T]],
) schnittstellen.FuncWriterFormat[schnittstellen.Set[T]] {
	return func(w io.Writer, e schnittstellen.Set[T]) (n int64, err error) {
		var n1 int

		if n1, err = io.WriteString(w, sf(e)); err != nil {
			n = int64(n1)
			err = errors.Wrap(err)
			return
		}

		n = int64(n1)

		return
	}
}
