package transaktion

import (
	"github.com/friedenberg/zit/src/delta/ts"
	sku "github.com/friedenberg/zit/src/foxtrot/sku"
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
