package iter2

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
)

func AddPtrOrReplaceIfGreater[T any, TPtr interfaces.Ptr[T]](
	c interfaces.MutableSetPtrLike[T, TPtr],
	l interfaces.Lessor2[T, TPtr],
	b TPtr,
) (err error) {
	a, ok := c.GetPtr(c.KeyPtr(b))

	if !ok || l.LessPtr(a, b) {
		return c.AddPtr(b)
	}

	return
}

func Parallel[T any](
	c interfaces.SetLike[T],
	f interfaces.FuncIter[T],
) (err error) {
	eg := iter.MakeErrorWaitGroupParallel()

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
