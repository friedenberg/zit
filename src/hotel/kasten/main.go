package kasten

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type Transacted = objekte.Transacted[
	Akte,
	*Akte,
	kennung.Kasten,
	*kennung.Kasten,
]

type CheckedOut = objekte.CheckedOut[
	Akte,
	*Akte,
	kennung.Kasten,
	*kennung.Kasten,
]
