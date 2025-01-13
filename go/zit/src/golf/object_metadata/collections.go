package object_metadata

import "code.linenisgreat.com/zit/go/zit/src/delta/heap"

type (
	Heap = heap.Heap[Metadata, *Metadata]
)

func MakeHeap() *Heap {
	return heap.Make[Metadata, *Metadata](
		Equaler,
		Lessor,
		Resetter,
	)
}
