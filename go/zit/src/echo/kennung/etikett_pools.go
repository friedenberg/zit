package kennung

import (
	"sync"

	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/bravo/pool"
)

var EtikettResetter etikettResetter

type etikettResetter struct{}

func (etikettResetter) Reset(e *Etikett) {
	e.value = ""
}

func (etikettResetter) ResetWith(a, b *Etikett) {
	a.value = b.value
}

func init() {
	etikettMapPool = pool.MakeValue[map[string]Etikett](
		func() map[string]Etikett {
			return make(map[string]Etikett)
		},
		func(v map[string]Etikett) {
			for k := range v {
				delete(v, k)
			}
		},
	)
}

var (
	etikettPool     schnittstellen.Pool[Etikett, *Etikett]
	etikettPoolOnce sync.Once
)

func GetEtikettPool() schnittstellen.Pool[Etikett, *Etikett] {
	etikettPoolOnce.Do(
		func() {
			etikettPool = pool.MakePool[Etikett, *Etikett](
				nil,
				EtikettResetter.Reset,
			)
		},
	)

	return etikettPool
}

var (
	etikettPtrMapPool     schnittstellen.PoolValue[map[string]*Etikett]
	etikettPtrMapPoolOnce sync.Once
)

func GetEtikettMapPtrPool() schnittstellen.PoolValue[map[string]*Etikett] {
	etikettPtrMapPoolOnce.Do(
		func() {
			etikettPtrMapPool = pool.MakeValue[map[string]*Etikett](
				func() map[string]*Etikett {
					return make(map[string]*Etikett)
				},
				func(v map[string]*Etikett) {
					for k := range v {
						// etikettPool.Put(e)
						delete(v, k)
					}
				},
			)
		},
	)

	return etikettPtrMapPool
}

var etikettMapPool schnittstellen.PoolValue[map[string]Etikett]
