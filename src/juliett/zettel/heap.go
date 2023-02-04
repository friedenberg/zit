package zettel

import "github.com/friedenberg/zit/src/charlie/collections"

type HeapTransacted = collections.Heap[Transacted, *Transacted]

func MakeHeapTransacted() HeapTransacted {
	return collections.MakeHeap[Transacted,*Transacted]()
}
