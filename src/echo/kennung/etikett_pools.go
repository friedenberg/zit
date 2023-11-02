package kennung

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/pool"
)

var (
	etikettPool       schnittstellen.Pool[Etikett, *Etikett]
	etikettPtrMapPool schnittstellen.PoolValue[map[string]*Etikett]
	etikettMapPool    schnittstellen.PoolValue[map[string]Etikett]
	EtikettResetter   etikettResetter
)

type etikettResetter struct{}

func (etikettResetter) Reset(e *Etikett) {
	e.value = ""
}

func init() {
	etikettPool = pool.MakePool[Etikett, *Etikett](
		nil,
		EtikettResetter.Reset,
	)

	etikettMapPool = pool.MakePoolValue[map[string]Etikett](
		func() map[string]Etikett {
			return make(map[string]Etikett)
		},
		func(v map[string]Etikett) {
			for k := range v {
				delete(v, k)
			}
		},
	)

	etikettPtrMapPool = pool.MakePoolValue[map[string]*Etikett](
		func() map[string]*Etikett {
			return make(map[string]*Etikett)
		},
		func(v map[string]*Etikett) {
			for k, e := range v {
				etikettPool.Put(e)
				delete(v, k)
			}
		},
	)
}

func GetEtikettPool() schnittstellen.Pool[Etikett, *Etikett] {
	return etikettPool
}
