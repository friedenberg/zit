package bestandsaufnahme

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Sku = sku.Transacted[
	kennung.Tai,
	*kennung.Tai,
]

type Transacted = sku.Transacted[
	kennung.Tai,
	*kennung.Tai,
]
