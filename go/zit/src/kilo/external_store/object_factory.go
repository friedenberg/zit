package external_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type ObjectFactory struct {
	interfaces.PoolValue[sku.ExternalLike]
	interfaces.Resetter3[sku.ExternalLike]
}

func (of *ObjectFactory) SetDefaultsIfNecessary() {
	if of.Resetter3 == nil {
		of.Resetter3 = pool.BespokeResetter[sku.ExternalLike]{
			FuncReset: func(e sku.ExternalLike) {
				sku.TransactedResetter.Reset(e.GetSku())
			},
			FuncResetWith: func(dst, src sku.ExternalLike) {
				sku.TransactedResetter.ResetWith(dst.GetSku(), src.GetSku())
			},
		}
	}

	if of.PoolValue == nil {
		of.PoolValue = pool.Bespoke[sku.ExternalLike]{
			FuncGet: func() sku.ExternalLike {
				return sku.GetTransactedPool().Get()
			},
			FuncPut: func(e sku.ExternalLike) {
				sku.GetTransactedPool().Put(e.(*sku.Transacted))
			},
		}
	}
}
