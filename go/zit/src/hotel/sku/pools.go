package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var (
	poolTransacted interfaces.Pool[Transacted, *Transacted]
	poolExternal   interfaces.Pool[External, *External]
	poolCheckedOut interfaces.Pool[CheckedOut, *CheckedOut]
)

func init() {
	poolTransacted = pool.MakePool(
		nil,
		TransactedResetter.Reset,
	)

	poolExternal = pool.MakePool(
		nil,
		TransactedResetter.Reset,
	)

	poolCheckedOut = pool.MakePool(
		nil,
		CheckedOutResetter.Reset,
	)
}

func GetTransactedPool() interfaces.Pool[Transacted, *Transacted] {
	return poolTransacted
}

func GetExternalPool() interfaces.Pool[External, *External] {
	return poolExternal
}

func GetCheckedOutPool() interfaces.Pool[CheckedOut, *CheckedOut] {
	return poolCheckedOut
}
