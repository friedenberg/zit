package zettel_transacted

import (
	"github.com/friedenberg/zit/src/delta/ts"
	"github.com/friedenberg/zit/zettel_named"
)

type Transacted struct {
	Named                 zettel_named.Named
	Kopf, Mutter, Schwanz ts.Time
}
