package kasten

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type Transacted = objekte.Transacted[
	Objekte,
	*Objekte,
	kennung.Kasten,
	*kennung.Kasten,
	Verzeichnisse,
	*Verzeichnisse,
]