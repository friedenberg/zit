package collections

import (
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
