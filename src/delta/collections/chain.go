package collections

import "github.com/friedenberg/zit/src/alfa/errors"

func MakeChain[T any](wfs ...WriterFunc[T]) WriterFunc[T] {
	return func(e T) (err error) {
		for _, w := range wfs {
			err = w(e)

			switch {
			case err == nil:
				continue

			case errors.IsEOF(err):
				err = nil
				return

			default:
				err = errors.Wrap(err)
				return
			}
		}

		return
	}
}
