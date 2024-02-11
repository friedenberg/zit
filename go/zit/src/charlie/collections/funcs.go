package collections

import (
	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
)

// TODO-P3 move to iter
func MakeWriterNoop[T any]() schnittstellen.FuncIter[T] {
	return func(e T) (err error) {
		return
	}
}

// TODO-P3 move to iter
func MakeTryFinally[T any](
	try schnittstellen.FuncIter[T],
	finally schnittstellen.FuncIter[T],
) schnittstellen.FuncIter[T] {
	return func(e T) (err error) {
		defer func() {
			err1 := finally(e)

			if err != nil {
				err = errors.MakeMulti(err, err1)
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
