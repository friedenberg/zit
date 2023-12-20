package heap

import (
	"container/heap"
	"sort"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/pool"
)

type Element interface{}

type ElementPtr[T Element] interface {
	schnittstellen.Ptr[T]
}

type heapPrivate[T Element, TPtr ElementPtr[T]] struct {
	Lessor     schnittstellen.Lessor3[TPtr]
	Resetter   schnittstellen.Resetter2[T, TPtr]
	Elements   []TPtr
	lastPopped TPtr
	p          schnittstellen.Pool[T, TPtr]
	equaler    schnittstellen.Equaler1[TPtr]
}

func (h *heapPrivate[T, TPtr]) GetPool() schnittstellen.Pool[T, TPtr] {
	if h.p == nil {
		h.p = pool.MakeFakePool[T, TPtr]()
	}

	return h.p
}

func (h heapPrivate[T, TPtr]) Len() int {
	return len(h.Elements)
}

func (h heapPrivate[T, TPtr]) Less(i, j int) (ok bool) {
	ok = h.Lessor.Less(h.Elements[i], h.Elements[j])
	return
}

func (h heapPrivate[T, TPtr]) Swap(i, j int) {
	h.Elements[i], h.Elements[j] = h.Elements[j], h.Elements[i]
}

func (h *heapPrivate[T, TPtr]) Push(x any) {
	h.Elements = append(h.Elements, x.(TPtr))
}

func (h *heapPrivate[T, TPtr]) discardDupes() {
	for h.lastPopped != nil &&
		h.Len() > 0 &&
		h.equaler.Equals(h.lastPopped, h.Elements[0]) {
		heap.Pop(h)
		// d := heap.Pop(h)
		// log.Debug().Printf("discarded: %s", d)
	}
}

func (h *heapPrivate[T, TPtr]) Pop() any {
	old := h.Elements
	n := len(old)
	x := old[n-1]
	h.Elements = old[0 : n-1]

	return x
}

func (h *heapPrivate[T, TPtr]) saveLastPopped(e TPtr) {
	if h.lastPopped == nil {
		h.lastPopped = h.GetPool().Get()
	}

	h.Resetter.ResetWith(h.lastPopped, e)
}

func (a heapPrivate[T, TPtr]) Copy() (b heapPrivate[T, TPtr]) {
	b = heapPrivate[T, TPtr]{
		Lessor:   a.Lessor,
		equaler:  a.equaler,
		p:        a.p,
		Resetter: a.Resetter,
		Elements: make([]TPtr, a.Len()),
	}

	copy(b.Elements, a.Elements)

	return
}

func (a heapPrivate[T, TPtr]) Sorted() (b heapPrivate[T, TPtr]) {
	b = a.Copy()
	sort.Sort(b)
	return
}
