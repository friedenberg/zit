package transaktion

import (
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/delta/objekte"
)

type Transaktion struct {
	ts.Time
	objekte.MutableSet
}

func MakeTransaktion(t ts.Time) Transaktion {
	return Transaktion{
		Time:       t,
		MutableSet: objekte.MakeMutableSet(),
	}
}
