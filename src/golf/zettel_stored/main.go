package zettel_stored

import (
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/objekte_stored"
)

type Stored = objekte_stored.Stored[zettel.Zettel, *zettel.Zettel]
