package transaktion

import (
	"github.com/friedenberg/zit/src/foxtrot/ts"
	"github.com/friedenberg/zit/src/golf/sku"
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
