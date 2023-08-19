package erworben

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type Transacted = sku.Transacted[
	kennung.Konfig,
	*kennung.Konfig,
]
