package zettel_transacted

import (
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/int_value"
	"github.com/friedenberg/zit/src/delta/objekte"
)

type Zettel struct {
	Named                 zettel_named.Zettel
	Kopf, Mutter, Schwanz ts.Time
	TransaktionIndex      int_value.IntValue
}

func (zt Zettel) IsNew() bool {
	return zt.Kopf == zt.Schwanz && zt.TransaktionIndex == 0
}

func (zt Zettel) ObjekteTransacted() objekte.ObjekteTransacted {
	return objekte.ObjekteTransacted{
		ObjekteWithIndex: objekte.ObjekteWithIndex{
			Objekte: objekte.Objekte{
				Gattung: gattung.Zettel,
				//TODO add Mutter
				Id:  &zt.Named.Hinweis,
				Sha: zt.Named.Stored.Sha,
			},
			Index: zt.TransaktionIndex,
		},
		Schwanz: zt.Schwanz,
	}
}

func (zt *Zettel) Reset() {
	zt.Kopf = ts.Time{}
	zt.Mutter = ts.Time{}
	zt.Schwanz = ts.Time{}
	zt.TransaktionIndex.Reset()

	zt.Named.Reset()
}
