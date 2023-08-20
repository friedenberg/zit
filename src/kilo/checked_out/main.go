package checked_out

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/juliett/objekte"
)

type (
	Etikett = objekte.CheckedOut[kennung.Etikett, *kennung.Etikett]
	Kasten  = objekte.CheckedOut[kennung.Kasten, *kennung.Kasten]
	Typ     = objekte.CheckedOut[kennung.Typ, *kennung.Typ]
	Zettel  = objekte.CheckedOut[kennung.Hinweis, *kennung.Hinweis]
)
