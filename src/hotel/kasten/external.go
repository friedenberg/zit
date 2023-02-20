package kasten

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type ExternalKeyer = objekte.ExternalKeyer[
	Objekte,
	*Objekte,
	kennung.Kasten,
	*kennung.Kasten,
]

type External = objekte.External[
	Objekte,
	*Objekte,
	kennung.Kasten,
	*kennung.Kasten,
]
