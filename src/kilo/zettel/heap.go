package zettel

import "github.com/friedenberg/zit/src/charlie/collections"

type HeapTransacted = collections.Heap[Transacted]

func MakeHeapTransacted() HeapTransacted {
	return collections.MakeHeap[Transacted]()
}
