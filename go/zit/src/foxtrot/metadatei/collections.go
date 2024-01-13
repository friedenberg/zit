package metadatei

import "github.com/friedenberg/zit/src/delta/heap"

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
