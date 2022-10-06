package verzeichnisse

import "sync"

type ZettelPool struct {
	inner *sync.Pool
}

func MakeZettelPool() ZettelPool {
	return ZettelPool{
		inner: &sync.Pool{
			New: func() interface{} {
				return &Zettel{}
			},
		},
	}
}

func (ip ZettelPool) Get() *Zettel {
	return ip.inner.Get().(*Zettel)
}

func (ip ZettelPool) Put(i *Zettel) {
	i.Reset()
	ip.inner.Put(i)
}
