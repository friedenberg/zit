package collections

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
)

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

func Multiplex[T any](
	e WriterFunc[T],
	producers ...func(WriterFunc[T]) error,
) (err error) {
	ch := make(chan error, len(producers))
	wg := &sync.WaitGroup{}
	wg.Add(len(producers))

	for _, p := range producers {
		go func(p func(WriterFunc[T]) error, ch chan<- error) {
			var err error

			defer func() {
				ch <- err
				wg.Done()
			}()

			if err = p(e); err != nil {
				err = errors.Wrap(err)
				return
			}
		}(p, ch)
	}

	wg.Wait()
	close(ch)

	for e := range ch {
		if e != nil {
			err = errors.MakeMulti(err, e)
		}
	}

	return
}
