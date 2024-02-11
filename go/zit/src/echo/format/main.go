package format

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

func MakeWriter[T any](
	wff schnittstellen.FuncWriterFormat[T],
	e T,
) schnittstellen.FuncWriter {
	return func(w io.Writer) (int64, error) {
		return wff(w, e)
	}
}

func MakeWriterOr[A schnittstellen.Stringer, B schnittstellen.Stringer](
	wffA schnittstellen.FuncWriterFormat[A],
	eA A,
	wffB schnittstellen.FuncWriterFormat[B],
	eB B,
) schnittstellen.FuncWriter {
	return func(w io.Writer) (int64, error) {
		if eA.String() == "" {
			return wffB(w, eB)
		} else {
			return wffA(w, eA)
		}
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
	sf schnittstellen.FuncString[schnittstellen.SetLike[T]],
) schnittstellen.FuncWriterFormat[schnittstellen.SetLike[T]] {
	return func(w io.Writer, e schnittstellen.SetLike[T]) (n int64, err error) {
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
