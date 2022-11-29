package objekte

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/echo/sku"
)

type Transacted[
	T gattung.Objekte,
	T1 gattung.ObjektePtr[T],
	T2 gattung.Identifier[T2],
	T3 gattung.IdentifierPtr[T2],
] struct {
	Named            Named[T, T1, T2, T3]
	Kopf, Schwanz    ts.Time
	Mutter           sku.Mutter
	TransaktionIndex int_value.IntValue
}

func (zt Transacted[T, T1, T2, T3]) IsNew() bool {
	return zt.Kopf == zt.Schwanz && zt.TransaktionIndex == 0
}

func (zt Transacted[T, T1, T2, T3]) SkuTransacted() sku.Transacted {
	return sku.Transacted{
		Indexed: sku.Indexed{
			Sku: sku.Sku{
				Gattung: zt.Named.Stored.Objekte.Gattung(),
				Mutter:  zt.Mutter,
				Id:      T3(&zt.Named.Kennung),
				Sha:     zt.Named.Stored.Sha,
			},
			Index: zt.TransaktionIndex,
		},
		Schwanz: zt.Schwanz,
	}
}

func (zt *Transacted[T, T1, T2, T3]) Reset() {
	zt.Kopf = ts.Time{}
	zt.Mutter[0] = ts.Time{}
	zt.Schwanz = ts.Time{}
	zt.TransaktionIndex.Reset()

	zt.Named.Reset()
}
