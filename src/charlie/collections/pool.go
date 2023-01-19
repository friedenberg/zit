package collections

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Resettable[T any] interface {
	Reset(T)
}

// TODO-P4 switch to interface
type Pool[T any] struct {
	inner *sync.Pool
}

func MakePool[T any]() *Pool[T] {
	return &Pool[T]{
		inner: &sync.Pool{
			New: func() interface{} {
				o := new(T)

				ii := interface{}(o)

				if r, ok := ii.(Resettable[*T]); ok {
					r.Reset(nil)
				}

				return o
			},
		},
	}
}

func (p Pool[T]) Apply(f WriterFunc[T], e T) (err error) {
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

func (ip Pool[T]) Get() *T {
	return ip.inner.Get().(*T)
}

func (ip Pool[T]) Put(i *T) (err error) {
	if i == nil {
		return
	}

	ii := interface{}(i)

	if r, ok := ii.(Resettable[*T]); ok {
		r.Reset(nil)
	}

	ip.inner.Put(i)

	return
}
