package zettel_transacted

import (
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
)

type Zettel struct {
	Named                 zettel_named.Zettel
	Kopf, Mutter, Schwanz ts.Time
}

func (zt Zettel) IsNew() bool {
	return zt.Kopf == zt.Schwanz
}

func (zt *Zettel) Reset() {
	zt.Kopf = ts.Time{}
	zt.Mutter = ts.Time{}
	zt.Schwanz = ts.Time{}

	zt.Named.Reset()
}
