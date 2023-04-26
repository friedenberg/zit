package typ

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type Transacted = objekte.Transacted[
	Akte,
	*Akte,
	kennung.Typ,
	*kennung.Typ,
	objekte.NilVerzeichnisse[Akte],
	*objekte.NilVerzeichnisse[Akte],
]

type ExternalKeyer = objekte.ExternalKeyer[
	Akte,
	*Akte,
	kennung.Typ,
	*kennung.Typ,
]

type External = objekte.External[
	Akte,
	*Akte,
	kennung.Typ,
	*kennung.Typ,
]

type CheckedOut = objekte.CheckedOut[
	Akte,
	*Akte,
	kennung.Typ,
	*kennung.Typ,
	objekte.NilVerzeichnisse[Akte],
	*objekte.NilVerzeichnisse[Akte],
]
