package collections

import "github.com/friedenberg/zit/src/alfa/errors"

type PooledChain[T any] []WriterFunc[*T]

func (pc PooledChain[T]) WriterWithPool(p Pool[T]) WriterFunc[*T] {
	return func(e *T) (err error) {
		return pc.Do(p, e)
	}
}

func (pc PooledChain[T]) Do(p Pool[T], e *T) (err error) {
	for _, w := range pc {
		err = w(e)

		switch {
		case e == nil && err == nil:
			return

		case err == nil:
			continue

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
