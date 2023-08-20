package etikett

import (
	"github.com/friedenberg/zit/src/delta/etikett_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type ExternalKeyer = objekte.ExternalKeyer[
	etikett_akte.V0,
	*etikett_akte.V0,
	kennung.Etikett,
	*kennung.Etikett,
]

type CheckedOut = objekte.CheckedOut[
	kennung.Etikett,
	*kennung.Etikett,
]
