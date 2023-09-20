package iter2

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
)

func AddPtrOrReplaceIfGreater[T any, TPtr schnittstellen.Ptr[T]](
	c schnittstellen.MutableSetPtrLike[T, TPtr],
	l schnittstellen.Lessor2[T, TPtr],
	b TPtr,
) (err error) {
	a, ok := c.GetPtr(c.KeyPtr(b))

	if !ok || l.LessPtr(a, b) {
		return c.AddPtr(b)
	}

	return
}

func Parallel[T any](
	c schnittstellen.SetLike[T],
	f schnittstellen.FuncIter[T],
) (err error) {
	eg := iter.MakeErrorWaitGroup()

	if err = c.Each(
		func(e T) (err error) {
			if !eg.Do(
				func() (err error) {
					return f(e)
				},
			) {
				err = iter.MakeErrStopIteration()
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = eg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
