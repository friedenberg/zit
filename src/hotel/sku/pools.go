package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/pool"
)

var (
	poolTransacted schnittstellen.Pool[Transacted, *Transacted]
	poolExternal   schnittstellen.Pool[External, *External]
)

func init() {
	poolTransacted = pool.MakePool[Transacted, *Transacted](
		nil,
		TransactedReseter.Reset,
	)

	poolExternal = pool.MakePool[External, *External](
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
