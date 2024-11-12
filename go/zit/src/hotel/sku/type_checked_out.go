package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
)

func makeCheckedOut() *CheckedOut {
	dst := GetCheckedOutPool().Get()
	return dst
}

func cloneFromTransactedCheckedOut(
	src *Transacted,
	newState checked_out_state.State,
) *CheckedOut {
	dst := GetCheckedOutPool().Get()
	TransactedResetter.ResetWith(dst.GetSku(), src)
	TransactedResetter.ResetWith(&dst.External, src)
	dst.State = newState
	return dst
}

func cloneCheckedOut(co *CheckedOut) *CheckedOut {
	return co.Clone()
}

type objectFactoryCheckedOut struct {
	interfaces.PoolValue[*CheckedOut]
	interfaces.Resetter3[*CheckedOut]
}

func (of *objectFactoryCheckedOut) SetDefaultsIfNecessary() objectFactoryCheckedOut {
	if of.Resetter3 == nil {
		of.Resetter3 = pool.BespokeResetter[*CheckedOut]{
			FuncReset: func(e *CheckedOut) {
				CheckedOutResetter.Reset(e)
			},
			FuncResetWith: func(dst, src *CheckedOut) {
				CheckedOutResetter.ResetWith(dst, src)
			},
		}
	}

	if of.PoolValue == nil {
		of.PoolValue = pool.Bespoke[*CheckedOut]{
			FuncGet: func() *CheckedOut {
				return GetCheckedOutPool().Get()
			},
			FuncPut: func(e *CheckedOut) {
				GetCheckedOutPool().Put(e)
			},
		}
	}

	return *of
}
