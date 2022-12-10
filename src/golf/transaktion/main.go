package transaktion

import (
	"github.com/friedenberg/zit/src/echo/sku"
	"github.com/friedenberg/zit/src/echo/ts"
)

type Transaktion struct {
	ts.Time
	sku.MutableSet
}

func MakeTransaktion(t ts.Time) Transaktion {
	return Transaktion{
		Time:       t,
		MutableSet: sku.MakeMutableSet(),
	}
}
