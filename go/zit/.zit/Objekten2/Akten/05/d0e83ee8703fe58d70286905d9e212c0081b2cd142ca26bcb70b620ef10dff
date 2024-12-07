package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var (
	poolTransacted interfaces.Pool[Transacted, *Transacted]
	poolCheckedOut interfaces.Pool[CheckedOut, *CheckedOut]
)

func init() {
	poolTransacted = pool.MakePool(
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

func GetCheckedOutPool() interfaces.Pool[CheckedOut, *CheckedOut] {
	return poolCheckedOut
}
