package erworben

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
)

type Transacted = sku.Transacted[
	kennung.Konfig,
	*kennung.Konfig,
]
