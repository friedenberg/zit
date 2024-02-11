package bestandsaufnahme

import (
	"code.linenisgreat.com/zit-go/src/charlie/sha"
	"code.linenisgreat.com/zit-go/src/delta/heap"
	"code.linenisgreat.com/zit-go/src/golf/ennui"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
)

type (
	Sku struct {
		sku.Transacted
		sha.Sha
		ennui.Range
	}

	SkuHeap = heap.Heap[Sku, *Sku]
)
