package sku

import (
	"github.com/friedenberg/zit/src/delta/collections"
)

type Sku2Heap = collections.Heap[Sku2]

func MakeSku2Heap() Sku2Heap {
	return collections.MakeHeap[Sku2]()
}
