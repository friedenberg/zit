package collections

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

// TODO-P4 switch to interface
type Pool[T any, T1 schnittstellen.Resetable[T]] struct {
	inner *sync.Pool
}

func MakePool[T any, T1 schnittstellen.Resetable[T]]() *Pool[T, T1] {
	return &Pool[T, T1]{
		inner: &sync.Pool{
			New: func() interface{} {
				o := new(T)
				T1(o).Reset()

				return o
			},
		},
	}
}

func (p Pool[T, T1]) Apply(f WriterFunc[T1], e T1) (err error) {
	err = f(e)

	switch {

	case IsDoNotRepool(err):
		err = nil
		return

	case IsStopIteration(err):
		err = nil
		p.Put(e)

	case err != nil:
		err = errors.Wrap(err)

		fallthrough

	default:
		p.Put(e)
	}

	return
}

func (ip Pool[T, T1]) Get() T1 {
	return ip.inner.Get().(T1)
}

func (ip Pool[T, T1]) Put(i T1) (err error) {
	if i == nil {
		return
	}

	i.Reset()
	ip.inner.Put(i)

	return
}
