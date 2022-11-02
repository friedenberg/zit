package collections

import "github.com/friedenberg/zit/src/alfa/errors"

type PooledChain[T any] []WriterFunc[*T]

func MakePooledChain[T any](wfs ...WriterFunc[*T]) PooledChain[T] {
	return PooledChain[T](wfs)
}

func (pc PooledChain[T]) WriterWithPool(p PoolLike[T]) WriterFunc[*T] {
	return func(e *T) (err error) {
		return pc.Do(p, e)
	}
}

func (pc PooledChain[T]) Do(p PoolLike[T], e *T) (err error) {
	for _, w := range pc {
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
