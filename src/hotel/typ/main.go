package typ

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type Transacted = objekte.Transacted[
	Objekte,
	*Objekte,
	kennung.Typ,
	*kennung.Typ,
	objekte.NilVerzeichnisse[Objekte],
	*objekte.NilVerzeichnisse[Objekte],
]

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

type CheckedOut = objekte.CheckedOut[
	Objekte,
	*Objekte,
	kennung.Typ,
	*kennung.Typ,
	objekte.NilVerzeichnisse[Objekte],
	*objekte.NilVerzeichnisse[Objekte],
]
