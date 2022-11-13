package transaktion

import (
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
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
