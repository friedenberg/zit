package collections

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

// TODO-P4 migrate to single writerfunc
func MakePooledChain[T any](
	p PoolLike[T],
	wfs ...schnittstellen.FuncIter[*T],
) schnittstellen.FuncIter[*T] {
	return func(e *T) (err error) {
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
