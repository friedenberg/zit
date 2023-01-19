package transaktion

import (
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/golf/sku"
)

type Transaktion struct {
	Time ts.Time
	Skus sku.MutableSet
}

func MakeTransaktion(t ts.Time) Transaktion {
	return Transaktion{
		Time: t,
		Skus: sku.MakeMutableSet(),
	}
}
