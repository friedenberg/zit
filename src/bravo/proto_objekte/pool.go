package proto_objekte

import "sync"

type Swimmer[T any] interface {
	*T
	Reset(*T)
}

type Pool[T Swimmer[T]] struct {
	inner *sync.Pool
}

func MakePool[T Swimmer[T]]() *Pool[T] {
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

func (ip Pool[T]) Put(i *T) {
	if i == nil {
		return
	}

	(*T(i)).Reset(nil)
	ip.inner.Put(i)
}
