package collections

import "github.com/friedenberg/zit/src/alfa/errors"

// TODO-P4 migrate to single writerfunc
func MakePooledChain[T any](p PoolLike[T], wfs ...WriterFunc[*T]) WriterFunc[*T] {
	return func(e *T) (err error) {
		for _, w := range wfs {
			err = w(e)

			switch {
			case err == nil:
				continue

			case errors.Is(err, ErrDoNotRepool{}):
				err = nil
				return

			case errors.IsEOF(err):
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
