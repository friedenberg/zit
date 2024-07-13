package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var poolTransacted interfaces.Pool[Transacted, *Transacted]

func init() {
	poolTransacted = pool.MakePool(
		nil,
		TransactedResetter.Reset,
	)
}

func GetTransactedPool() interfaces.Pool[Transacted, *Transacted] {
	return poolTransacted
}
