package external_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func cloneCheckedOut(co *sku.CheckedOut) *sku.CheckedOut {
	return co.Clone()
}

type objectFactoryCheckedOut struct {
	interfaces.PoolValue[*sku.CheckedOut]
	interfaces.Resetter3[*sku.CheckedOut]
}

func (of *objectFactoryCheckedOut) SetDefaultsIfNecessary() objectFactoryCheckedOut {
	if of.Resetter3 == nil {
		of.Resetter3 = pool.BespokeResetter[*sku.CheckedOut]{
			FuncReset: func(e *sku.CheckedOut) {
				sku.CheckedOutResetter.Reset(e)
			},
			FuncResetWith: func(dst, src *sku.CheckedOut) {
				sku.CheckedOutResetter.ResetWith(dst, src)
			},
		}
	}

	if of.PoolValue == nil {
		of.PoolValue = pool.Bespoke[*sku.CheckedOut]{
			FuncGet: func() *sku.CheckedOut {
				return sku.GetCheckedOutPool().Get()
			},
			FuncPut: func(e *sku.CheckedOut) {
				sku.GetCheckedOutPool().Put(e)
			},
		}
	}

	return *of
}
