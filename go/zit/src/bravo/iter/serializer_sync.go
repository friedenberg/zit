package iter

import (
	"sync"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
)

func MakeSyncSerializer[T any](
	wf schnittstellen.FuncIter[T],
) schnittstellen.FuncIter[T] {
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
