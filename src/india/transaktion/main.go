package transaktion

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Transaktion struct {
	Time kennung.Time
	Skus sku.MutableSet
}

func MakeTransaktion(t kennung.Time) Transaktion {
	return Transaktion{
		Time: t,
		Skus: sku.MakeMutableSet(),
	}
}
