package zettel_stored

import (
	"github.com/friedenberg/zit/src/foxtrot/objekte"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type Stored = objekte.Stored[zettel.Zettel, *zettel.Zettel]
