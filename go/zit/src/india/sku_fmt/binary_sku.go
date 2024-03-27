package sku_fmt

import (
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/golf/ennui"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type SkuWithSigil struct {
	*sku.Transacted
	kennung.Sigil
}

type Sku struct {
	SkuWithSigil
	ennui.Range
}
