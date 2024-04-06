package heap

import (
	"sort"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

func Make[T Element, TPtr ElementPtr[T]](
	equaler schnittstellen.Equaler1[TPtr],
	lessor schnittstellen.Lessor3[TPtr],
	resetter schnittstellen.Resetter2[T, TPtr],
) *Heap[T, TPtr] {
	return &Heap[T, TPtr]{
		h: heapPrivate[T, TPtr]{
			equaler:  equaler,
			Lessor:   lessor,
			Resetter: resetter,
			Elements: make([]TPtr, 0),
		},
	}
}

func MakeHeapFromSliceUnsorted[T Element, TPtr ElementPtr[T]](
	equaler schnittstellen.Equaler1[TPtr],
	lessor schnittstellen.Lessor3[TPtr],
	resetter schnittstellen.Resetter2[T, TPtr],
	s []TPtr,
) *Heap[T, TPtr] {
	h := heapPrivate[T, TPtr]{
		Lessor:   lessor,
		Resetter: resetter,
		Elements: s,
		equaler:  equaler,
	}

	sort.Sort(h)

	return &Heap[T, TPtr]{
		h: h,
	}
}

func MakeHeapFromSlice[T Element, TPtr ElementPtr[T]](
	equaler schnittstellen.Equaler1[TPtr],
	lessor schnittstellen.Lessor3[TPtr],
	resetter schnittstellen.Resetter2[T, TPtr],
	s []TPtr,
) *Heap[T, TPtr] {
	h := heapPrivate[T, TPtr]{
		Lessor:   lessor,
		Resetter: resetter,
		Elements: s,
		equaler:  equaler,
	}

	return &Heap[T, TPtr]{
		h: h,
	}
}
