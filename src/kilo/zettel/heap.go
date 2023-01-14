package zettel

import "github.com/friedenberg/zit/src/delta/collections"

type HeapTransacted = collections.Heap[Transacted]

func MakeHeapTransacted() HeapTransacted {
	return collections.MakeHeap[Transacted]()
}
