package external_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// type SkuType = *sku.CheckedOut
type SkuType = sku.ExternalLike

type ObjectFactory struct {
	interfaces.PoolValue[SkuType]
	interfaces.Resetter3[SkuType]
}

func (of *ObjectFactory) SetDefaultsIfNecessary() {
	if of.Resetter3 == nil {
		of.Resetter3 = pool.BespokeResetter[SkuType]{
			FuncReset: func(e SkuType) {
				sku.TransactedResetter.Reset(e.GetSku())
			},
			FuncResetWith: func(dst, src SkuType) {
				sku.TransactedResetter.ResetWith(dst.GetSku(), src.GetSku())
			},
		}
	}

	if of.PoolValue == nil {
		of.PoolValue = pool.Bespoke[SkuType]{
			FuncGet: func() SkuType {
				return sku.GetTransactedPool().Get()
			},
			FuncPut: func(e SkuType) {
				sku.GetTransactedPool().Put(e.(*sku.Transacted))
			},
		}
	}
}
