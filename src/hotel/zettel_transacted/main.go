package zettel_transacted

import (
	"github.com/friedenberg/zit/src/charlie/ts"
	"github.com/friedenberg/zit/src/golf/zettel_named"
)

type Zettel struct {
	Named                 zettel_named.Zettel
	Kopf, Mutter, Schwanz ts.Time
}
