package zettel_transacted

import "sync"

type PoolGetter interface {
	ZettelTransactedPool() *Pool
}

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
	i.Reset()
	ip.inner.Put(i)
}

func (ip Pool) WriteZettelVerzeichnisse(z *Zettel) (err error) {
	ip.Put(z)
	return
}
