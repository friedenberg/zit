package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/pool"
)

var poolTransacted schnittstellen.Pool[Transacted, *Transacted]

func init() {
	poolTransacted = pool.MakePool(
		nil,
		TransactedResetter.Reset,
	)
}

func GetTransactedPool() schnittstellen.Pool[Transacted, *Transacted] {
	return poolTransacted
}
