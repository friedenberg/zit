package object_metadata

import "code.linenisgreat.com/zit/go/zit/src/delta/heap"

type (
	Heap = heap.Heap[Metadatei, *Metadatei]
)

func MakeHeap() *Heap {
	return heap.Make[Metadatei, *Metadatei](
		Equaler,
		Lessor,
		Resetter,
	)
}
