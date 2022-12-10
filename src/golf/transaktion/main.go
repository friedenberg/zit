package transaktion

import (
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/foxtrot/ts"
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
