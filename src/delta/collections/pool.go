package collections

import (
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Resettable[T any] interface {
	Reset(T)
}

type Pool[T any] struct {
	inner *sync.Pool
}

func MakePool[T any]() *Pool[T] {
	return &Pool[T]{
		inner: &sync.Pool{
			New: func() interface{} {
				return new(T)
			},
		},
	}
}

func (p Pool[T]) Apply(f WriterFunc[T], e T) (err error) {
	err = f(e)

	switch {

	case errors.Is(err, ErrDoNotRepool{}):
		err = nil
		return

	case errors.IsEOF(err):
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
