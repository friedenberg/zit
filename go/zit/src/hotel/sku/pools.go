package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var (
	poolTransacted interfaces.Pool[Transacted, *Transacted]
	poolExternal   interfaces.Pool[External, *External]
)

func init() {
	poolTransacted = pool.MakePool(
		nil,
		TransactedResetter.Reset,
	)

	poolExternal = pool.MakePool(
		nil,
		ExternalResetter.Reset,
	)
}

func GetTransactedPool() interfaces.Pool[Transacted, *Transacted] {
	return poolTransacted
}

func GetExternalPool() interfaces.Pool[External, *External] {
	return poolExternal
}
