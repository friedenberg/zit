package collections

import (
	"sync"
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

func (ip Pool[T]) Get() *T {
	return ip.inner.Get().(*T)
}

func (ip Pool[T]) Put(i *T) (err error) {
	// errors.Err().Print("tried to put")
	// return
	if i == nil {
		panic("tried to pool nil")
	}

	ii := interface{}(i)

	if r, ok := ii.(Resettable[*T]); ok {
		r.Reset(nil)
	}

	ip.inner.Put(i)

	return
}
