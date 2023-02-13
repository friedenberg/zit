package collections

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type pool[T any, TPtr schnittstellen.Resetable[T]] struct {
	inner *sync.Pool
}

func MakePool[T any, TPtr schnittstellen.Resetable[T]]() *pool[T, TPtr] {
	return &pool[T, TPtr]{
		inner: &sync.Pool{
			New: func() interface{} {
				o := new(T)
				TPtr(o).Reset()

				return o
			},
		},
	}
}

func (p pool[T, TPtr]) Apply(f schnittstellen.FuncIter[T], e T) (err error) {
	err = f(e)

	switch {

	case IsDoNotRepool(err):
		err = nil
		return

	case IsStopIteration(err):
		err = nil
		p.Put(&e)

	case err != nil:
		err = errors.Wrap(err)

		fallthrough

	default:
		p.Put(&e)
	}

	return
}

func (ip pool[T, TPtr]) Get() TPtr {
	return ip.inner.Get().(TPtr)
}

func (ip pool[T, TPtr]) Put(i TPtr) (err error) {
	if i == nil {
		return
	}

	i.Reset()
	ip.inner.Put(i)

	return
}
