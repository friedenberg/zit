package stream_index

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/object_probe_index"
)

type skuWithSigil struct {
	*sku.Transacted
	ids.Sigil
}

type skuWithRangeAndSigil struct {
	skuWithSigil
	object_probe_index.Range
}
