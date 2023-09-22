package sku

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/pool"
)

var (
	poolTransacted schnittstellen.Pool[Transacted2, *Transacted2]
	poolExternal   schnittstellen.Pool[External2, *External2]
)

func init() {
	poolTransacted = pool.MakePool[Transacted2, *Transacted2](
		nil,
		nil,
	)

	poolExternal = pool.MakePool[External2, *External2](
		nil,
		nil,
	)
}

func GetTransactedPool() schnittstellen.Pool[Transacted2, *Transacted2] {
	return poolTransacted
}

func GetExternalPool() schnittstellen.Pool[External2, *External2] {
	return poolExternal
}
