package bestandsaufnahme

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type Sku = sku.Transacted[
	kennung.Tai,
	*kennung.Tai,
]

type Transacted = objekte.Transacted[
	Akte,
	*Akte,
	kennung.Tai,
	*kennung.Tai,
]
