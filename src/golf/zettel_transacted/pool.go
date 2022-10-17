package zettel_transacted

import "sync"

type Pool struct {
	inner *sync.Pool
}

func MakePool() *Pool {
	return &Pool{
		inner: &sync.Pool{
			New: func() interface{} {
				return &Zettel{}
			},
		},
	}
}

func (ip Pool) Get() *Zettel {
	return ip.inner.Get().(*Zettel)
}

func (ip Pool) Put(i *Zettel) {
	if i == nil {
		return
	}

	i.Reset()
	ip.inner.Put(i)
}

func (ip Pool) WriteZettelTransacted(z *Zettel) (err error) {
	ip.Put(z)
	return
}
