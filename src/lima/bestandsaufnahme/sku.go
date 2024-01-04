package bestandsaufnahme

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/heap"
	"github.com/friedenberg/zit/src/golf/ennui"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type (
	Sku struct {
		sku.Transacted
		sha.Sha
		ennui.Range
	}

	SkuHeap = heap.Heap[Sku, *Sku]
)
