package kasten

import (
	"github.com/friedenberg/zit/src/delta/kasten_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type ExternalKeyer = objekte.ExternalKeyer[
	kasten_akte.V0,
	*kasten_akte.V0,
	kennung.Kasten,
	*kennung.Kasten,
]

type External = sku.External[
	kennung.Kasten,
	*kennung.Kasten,
]
