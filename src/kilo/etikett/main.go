package etikett

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type ExternalKeyer = objekte.ExternalKeyer[
	Akte,
	*Akte,
	kennung.Etikett,
	*kennung.Etikett,
]

type CheckedOut = objekte.CheckedOut[
	kennung.Etikett,
	*kennung.Etikett,
]
