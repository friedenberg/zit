package etikett

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type Transacted = sku.Transacted[
	kennung.Etikett,
	*kennung.Etikett,
]

type ExternalKeyer = objekte.ExternalKeyer[
	Akte,
	*Akte,
	kennung.Etikett,
	*kennung.Etikett,
]

type External = objekte.External[
	Akte,
	*Akte,
	kennung.Etikett,
	*kennung.Etikett,
]

type CheckedOut = objekte.CheckedOut[
	Akte,
	*Akte,
	kennung.Etikett,
	*kennung.Etikett,
]
