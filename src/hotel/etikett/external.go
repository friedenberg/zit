package etikett

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type ExternalKeyer = objekte.ExternalKeyer[
	Objekte,
	*Objekte,
	kennung.Etikett,
	*kennung.Etikett,
]

type External = objekte.External[
	Objekte,
	*Objekte,
	kennung.Etikett,
	*kennung.Etikett,
]
