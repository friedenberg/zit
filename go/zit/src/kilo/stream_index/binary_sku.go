package stream_index

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/india/object_probe_index"
	"code.linenisgreat.com/zit/go/zit/src/india/sku"
)

type skuWithSigil struct {
	*sku.Transacted
	ids.Sigil
}

type skuWithRangeAndSigil struct {
	skuWithSigil
	object_probe_index.Range
}
