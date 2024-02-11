package sku

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/bravo/pool"
)

var (
	poolTransacted schnittstellen.Pool[Transacted, *Transacted]
	poolExternal   schnittstellen.Pool[External, *External]
	poolCheckedOut schnittstellen.Pool[CheckedOut, *CheckedOut]
)

func init() {
	poolTransacted = pool.MakePool[Transacted, *Transacted](
		nil,
		TransactedResetter.Reset,
	)

	poolExternal = pool.MakePool[External, *External](
		nil,
		nil,
	)

	poolCheckedOut = pool.MakePool[CheckedOut, *CheckedOut](
		nil,
		nil,
	)
}

func GetTransactedPool() schnittstellen.Pool[Transacted, *Transacted] {
	return poolTransacted
}

func GetExternalPool() schnittstellen.Pool[External, *External] {
	return poolExternal
}

func GetCheckedOutPool() schnittstellen.Pool[CheckedOut, *CheckedOut] {
	return poolCheckedOut
}
