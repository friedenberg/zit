package zettel_stored

import (
	objekte "github.com/friedenberg/zit/src/delta/objekte"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type Stored = objekte.Stored[zettel.Zettel, *zettel.Zettel]
