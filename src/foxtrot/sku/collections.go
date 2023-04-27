package sku

import (
	"github.com/friedenberg/zit/src/charlie/collections"
)

type Sku2Heap = collections.Heap[Sku, *Sku]

func MakeSku2Heap() Sku2Heap {
	return collections.MakeHeap[Sku, *Sku]()
}
