package quiter

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func ErrorWaitGroupApply[T any](
	wg ErrorWaitGroup,
	s interfaces.SetLike[T],
	f interfaces.FuncIter[T],
) (d bool) {
	if err := s.Each(
		func(e T) (err error) {
			if !wg.Do(
				func() error {
					return f(e)
				},
			) {
				err = errors.MakeErrStopIteration()
			}

			return
		},
	); err != nil {
		d = true
	}

	return
}
