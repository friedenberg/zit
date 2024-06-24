package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var (
	poolTransacted schnittstellen.Pool[Transacted, *Transacted]
	poolExternal   schnittstellen.Pool[ExternalFS, *ExternalFS]
	poolCheckedOut schnittstellen.Pool[CheckedOutFS, *CheckedOutFS]
)

func init() {
	poolTransacted = pool.MakePool(
		nil,
		TransactedResetter.Reset,
	)

	poolExternal = pool.MakePool[ExternalFS](
		nil,
		nil,
	)

	poolCheckedOut = pool.MakePool[CheckedOutFS](
		nil,
		nil,
	)
}

func GetTransactedPool() schnittstellen.Pool[Transacted, *Transacted] {
	return poolTransacted
}

func GetExternalPool() schnittstellen.Pool[ExternalFS, *ExternalFS] {
	return poolExternal
}

func GetCheckedOutPool() schnittstellen.Pool[CheckedOutFS, *CheckedOutFS] {
	return poolCheckedOut
}
