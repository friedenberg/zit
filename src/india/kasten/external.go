package kasten

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type ExternalKeyer = objekte.ExternalKeyer[
	Akte,
	*Akte,
	kennung.Kasten,
	*kennung.Kasten,
]

type External = sku.External[
	kennung.Kasten,
	*kennung.Kasten,
]
