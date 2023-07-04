package kasten

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type ExternalKeyer = objekte.ExternalKeyer[
	Akte,
	*Akte,
	kennung.Kasten,
	*kennung.Kasten,
]

type External = objekte.External[
	Akte,
	*Akte,
	kennung.Kasten,
	*kennung.Kasten,
]
