package transaktion

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
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
