package objekte_stored

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/int_value"
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/foxtrot/objekte"
)

type Transacted[T any, T1 ObjektePtr[T], T2 Identifier[T2], T3 IdentifierPtr[T2]] struct {
	Named                 Named[T, T1, T2, T3]
	Kopf, Mutter, Schwanz ts.Time
	TransaktionIndex      int_value.IntValue
}

func (zt Transacted[T, T1, T2, T3]) IsNew() bool {
	return zt.Kopf == zt.Schwanz && zt.TransaktionIndex == 0
}

func (zt Transacted[T, T1, T2, T3]) ObjekteTransacted() objekte.ObjekteTransacted {
	return objekte.ObjekteTransacted{
		ObjekteWithIndex: objekte.ObjekteWithIndex{
			Objekte: objekte.Objekte{
				Gattung: gattung.Zettel,
				//TODO add Mutter
				Id:  T3(&zt.Named.Kennung),
				Sha: zt.Named.Stored.Sha,
			},
			Index: zt.TransaktionIndex,
		},
		Schwanz: zt.Schwanz,
	}
}

func (zt *Transacted[T, T1, T2, T3]) Reset() {
	zt.Kopf = ts.Time{}
	zt.Mutter = ts.Time{}
	zt.Schwanz = ts.Time{}
	zt.TransaktionIndex.Reset()

	zt.Named.Reset()
}
