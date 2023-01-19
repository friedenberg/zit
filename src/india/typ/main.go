package typ

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/objekte"
)

type Transacted = objekte.Transacted[
	Objekte,
	*Objekte,
	kennung.Typ,
	*kennung.Typ,
	objekte.NilVerzeichnisse[Objekte],
	*objekte.NilVerzeichnisse[Objekte],
]
