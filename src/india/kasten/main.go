package kasten

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type CheckedOut = objekte.CheckedOut[
	kennung.Kasten,
	*kennung.Kasten,
]