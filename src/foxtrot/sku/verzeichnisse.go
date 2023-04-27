package sku

import (
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/ts"
)

type Verzeichnisse struct {
	Mutter           Mutter
	TransactionIndex values.Int
	Kopf, Schwanz    ts.Time
}

func (v *Verzeichnisse) SetFromSku(s Sku) {
	v.Kopf = s.GetTime()
	v.Schwanz = s.GetTime()
}
