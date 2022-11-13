package collections

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func MakeWriterDoNotRepool[T any]() WriterFunc[*T] {
	return func(e *T) (err error) {
		err = ErrDoNotRepool{}
		return
	}
}

func MakeWriterNoop[T any]() WriterFunc[T] {
	return func(e T) (err error) {
		return
	}
}

func MakeTryFinally[T any](
	try WriterFunc[T],
	finally WriterFunc[T],
) WriterFunc[T] {
	return func(e T) (err error) {
		defer func() {
			err1 := finally(e)

			if err != nil {
				err = errors.MakeErrorMultiOrNil(err, err1)
			} else {
				err = err1
			}
		}()

		if err = try(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func MakeWriterFormatFunc[T any](
	wff WriterFuncFormat[T],
	e *T,
) Writer {
	return func(w io.Writer) (int64, error) {
		return wff(w, e)
	}
}

// TODO rename
func MakeWriterLiteral(
	f string,
	vs ...interface{},
) Writer {
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

func MakeWriterFormatStringer[T fmt.Stringer]() WriterFuncFormat[T] {
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
