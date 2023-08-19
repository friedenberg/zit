package transacted

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
)

type (
	Etikett = sku.Transacted[kennung.Etikett, *kennung.Etikett]
	Kasten  = sku.Transacted[kennung.Kasten, *kennung.Kasten]
	Konfig  = sku.Transacted[kennung.Konfig, *kennung.Konfig]
	Typ     = sku.Transacted[kennung.Typ, *kennung.Typ]
	Zettel  = sku.Transacted[kennung.Hinweis, *kennung.Hinweis]
)
