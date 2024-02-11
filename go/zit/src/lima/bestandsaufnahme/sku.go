package bestandsaufnahme

import (
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/heap"
	"code.linenisgreat.com/zit/src/golf/ennui"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type (
	Sku struct {
		sku.Transacted
		sha.Sha
		ennui.Range
	}

	SkuHeap = heap.Heap[Sku, *Sku]
)
