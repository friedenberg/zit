package etikett

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type Transacted = objekte.Transacted[
	Akte,
	*Akte,
	kennung.Etikett,
	*kennung.Etikett,
	objekte.NilVerzeichnisse[Akte],
	*objekte.NilVerzeichnisse[Akte],
]

type ExternalKeyer = objekte.ExternalKeyer[
	Akte,
	*Akte,
	kennung.Etikett,
	*kennung.Etikett,
]

type External = objekte.External[
	Akte,
	*Akte,
	kennung.Etikett,
	*kennung.Etikett,
]

type CheckedOut = objekte.CheckedOut[
	Akte,
	*Akte,
	kennung.Etikett,
	*kennung.Etikett,
	objekte.NilVerzeichnisse[Akte],
	*objekte.NilVerzeichnisse[Akte],
]
