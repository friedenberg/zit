package alfred

import (
	"sync"

	"code.linenisgreat.com/zit/src/charlie/catgut"
)

type ItemPool struct {
	inner *sync.Pool
}

func MakeItemPool() ItemPool {
	return ItemPool{
		inner: &sync.Pool{
			New: func() interface{} {
				return &Item{
					Match: &catgut.String{},
				}
			},
		},
	}
}

func (ip ItemPool) Get() *Item {
	return ip.inner.Get().(*Item)
}

func (ip ItemPool) Put(i *Item) {
	if i == nil {
		panic("item was nil")
	}

	i.Reset()
	ip.inner.Put(i)
}
