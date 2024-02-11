package metadatei

import "code.linenisgreat.com/zit/src/delta/heap"

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
