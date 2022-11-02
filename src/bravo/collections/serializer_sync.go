package collections

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type SyncSerializer[T any] struct {
	wf WriterFunc[T]
	l  sync.Locker
}

func MakeSyncSerializer[T any](wf WriterFunc[T]) SyncSerializer[T] {
	return SyncSerializer[T]{
		wf: wf,
		l:  &sync.Mutex{},
	}
}

func (s SyncSerializer[T]) Do(e T) (err error) {
	s.l.Lock()
	defer s.l.Unlock()

	if err = s.wf(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
