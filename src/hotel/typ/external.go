package typ

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type ExternalKeyer = objekte.ExternalKeyer[
	Objekte,
	*Objekte,
	kennung.Typ,
	*kennung.Typ,
]

type External = objekte.External[
	Objekte,
	*Objekte,
	kennung.Typ,
	*kennung.Typ,
]
