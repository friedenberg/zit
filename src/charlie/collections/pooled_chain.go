package collections

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func MakePooledChain[T any, TPtr schnittstellen.Ptr[T]](
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

			case IsStopIteration(err):
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
