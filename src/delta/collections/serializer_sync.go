package collections

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func MakeSyncSerializer[T any](wf WriterFunc[T]) WriterFunc[T] {
	l := &sync.Mutex{}

	return func(e T) (err error) {
		l.Lock()
		defer l.Unlock()

		if err = wf(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
