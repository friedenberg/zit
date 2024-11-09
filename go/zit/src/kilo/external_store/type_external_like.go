package external_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func makeExternalLike() sku.ExternalLike {
	dst := sku.GetTransactedPool().Get()
	return dst
}

func cloneFromTransactedExternalLike(src *sku.Transacted) sku.ExternalLike {
	dst := sku.GetTransactedPool().Get()
	sku.TransactedResetter.ResetWith(dst, src)
	return dst
}

func cloneExternalLike(el sku.ExternalLike) sku.ExternalLike {
	return el.CloneExternalLike()
}

type objectFactoryExternalLike struct {
	interfaces.PoolValue[sku.ExternalLike]
	interfaces.Resetter3[sku.ExternalLike]
}

func (of *objectFactoryExternalLike) SetDefaultsIfNecessary() objectFactoryExternalLike {
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

	return *of
}
