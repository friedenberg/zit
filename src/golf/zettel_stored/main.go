package zettel_stored

import (
	"github.com/friedenberg/zit/src/delta/objekte_stored"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type Stored = objekte_stored.Stored[zettel.Zettel, *zettel.Zettel]
