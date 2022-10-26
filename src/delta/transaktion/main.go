package transaktion

import (
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/objekte"
)

type Transaktion struct {
	ts.Time
	Objekten []objekte.Objekte
}

func (t *Transaktion) AddObjekte(o objekte.Objekte) {
	t.Objekten = append(t.Objekten, o)
}
