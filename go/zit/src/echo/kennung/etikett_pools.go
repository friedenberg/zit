package kennung

import (
	"sync"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/pool"
)

var (
	etikettPool     schnittstellen.Pool[Etikett, *Etikett]
	etikettPoolOnce sync.Once
)

func init() {
	etikettMapPool = pool.MakeValue(
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

type etikettResetter struct{}

func (etikettResetter) Reset(e *etikett) {
	e.value = ""
	e.virtual = false
	e.dependentLeaf = false
}

func (etikettResetter) ResetWith(a, b *etikett) {
	a.value = b.value
	a.virtual = b.virtual
	a.dependentLeaf = b.dependentLeaf
}

type etikett2Resetter struct{}

func (etikett2Resetter) Reset(e *etikett2) {
	e.value.Reset()
	e.virtual = false
	e.dependentLeaf = false
}

func (etikett2Resetter) ResetWith(a, b *etikett2) {
	b.value.CopyTo(a.value)
	a.virtual = b.virtual
	a.dependentLeaf = b.dependentLeaf
}

var (
	etikettPtrMapPool     schnittstellen.PoolValue[map[string]*Etikett]
	etikettPtrMapPoolOnce sync.Once
)

func GetEtikettPool() schnittstellen.Pool[Etikett, *Etikett] {
	etikettPoolOnce.Do(
		func() {
			etikettPool = pool.MakePool(
				func() *Etikett {
					e := &Etikett{}
					e.init()
					return e
				},
				EtikettResetter.Reset,
			)
		},
	)

	return etikettPool
}

func GetEtikettMapPtrPool() schnittstellen.PoolValue[map[string]*Etikett] {
	etikettPtrMapPoolOnce.Do(
		func() {
			etikettPtrMapPool = pool.MakeValue(
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
