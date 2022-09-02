package zettel_transacted

import (
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/src/hotel/zettel_named"
)

type Zettel struct {
	Named                 zettel_named.Zettel
	Kopf, Mutter, Schwanz ts.Time
}
