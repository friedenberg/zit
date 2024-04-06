package store_verzeichnisse

import (
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/golf/ennui"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type skuWithSigil struct {
	*sku.Transacted
	kennung.Sigil
}

type skuWithRangeAndSigil struct {
	skuWithSigil
	ennui.Range
}
