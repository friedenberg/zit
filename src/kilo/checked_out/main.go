package checked_out

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type (
	Typ     = objekte.CheckedOut[kennung.Typ, *kennung.Typ]
	Etikett = objekte.CheckedOut[kennung.Etikett, *kennung.Etikett]
)
