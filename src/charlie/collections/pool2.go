package collections

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type pool2[T any, TPtr schnittstellen.Resetable[T]] struct {
	inner *sync.Pool
}

func MakePool2[T any, TPtr schnittstellen.Resetable[T]]() *pool2[T, TPtr] {
	return &pool2[T, TPtr]{
		inner: &sync.Pool{
			New: func() interface{} {
				o := new(T)
				TPtr(o).Reset()

				return o
			},
		},
	}
}

func (p pool2[T, TPtr]) Apply(f WriterFunc[T], e T) (err error) {
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

func (ip pool2[T, TPtr]) Get() TPtr {
	return ip.inner.Get().(TPtr)
}

func (ip pool2[T, TPtr]) Put(i TPtr) (err error) {
	if i == nil {
		return
	}

	i.Reset()
	ip.inner.Put(i)

	return
}