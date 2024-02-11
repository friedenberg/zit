package pool

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

func MakeWriterDoNotRepool[T any]() schnittstellen.FuncIter[*T] {
	return func(e *T) (err error) {
		err = ErrDoNotRepool
		return
	}
}

func MakePooledChain[T schnittstellen.Poolable[T], TPtr schnittstellen.PoolablePtr[T]](
	p schnittstellen.Pool[T, TPtr],
	wfs ...schnittstellen.FuncIter[TPtr],
) schnittstellen.FuncIter[TPtr] {
	return func(e TPtr) (err error) {
		for _, w := range wfs {
			err = w(e)

			switch {
			case err == nil:
				continue

			case IsDoNotRepool(err):
				err = nil
				return

			case errors.IsStopIteration(err):
				err = nil
				p.Put(e)
				return

			default:
				p.Put(e)
				err = errors.Wrap(err)
				return
			}
		}

		p.Put(e)

		return
	}
}
