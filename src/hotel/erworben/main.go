package erworben

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type Transacted = objekte.Transacted[
	Akte,
	*Akte,
	kennung.Konfig,
	*kennung.Konfig,
	objekte.NilVerzeichnisse[Akte],
	*objekte.NilVerzeichnisse[Akte],
]
