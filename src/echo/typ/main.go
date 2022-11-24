package typ

import (
	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/delta/objekte"
)

type Stored = objekte.Stored[Akte, *Akte]
type Named = objekte.Named[Akte, *Akte, kennung.Typ, *kennung.Typ]
