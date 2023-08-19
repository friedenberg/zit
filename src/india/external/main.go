package external

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type (
	Etikett = sku.External[kennung.Etikett, *kennung.Etikett]
	Kasten  = sku.External[kennung.Kasten, *kennung.Kasten]
	Konfig  = sku.External[kennung.Konfig, *kennung.Konfig]
	Typ     = sku.External[kennung.Typ, *kennung.Typ]
	Zettel  = sku.External[kennung.Hinweis, *kennung.Hinweis]
)
