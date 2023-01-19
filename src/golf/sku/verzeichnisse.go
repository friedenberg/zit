package sku

import (
	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/echo/ts"
)

type Verzeichnisse struct {
	Mutter           Mutter
	TransactionIndex int_value.IntValue
	Kopf, Schwanz    ts.Time
}

func (v *Verzeichnisse) SetFromSku(s Sku) {
	v.Kopf = s.GetTime()
	v.Schwanz = s.GetTime()
}

func (v *Verzeichnisse) SetFromSku2(s Sku2) {
	v.Kopf = s.GetTime()
	v.Schwanz = s.GetTime()
}
