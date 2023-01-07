package sku

import (
	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/foxtrot/ts"
)

type Verzeichnisse struct {
	Mutter     Mutter
	TransactionIndex int_value.IntValue
	Kopf, Schwanz    ts.Time
}
