package collections

import (
	"container/heap"
	"sort"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
)

type Lessor[T any] interface {
	Less(T) bool
}

type Equaler[T any] interface {
	Equals(*T) bool
}

type HeapElement[T gattung.Element] interface {
	Equaler[T]
	Lessor[T]
}

type HeapElementPtr[T gattung.Element] interface {
	gattung.ElementPtr[T]
	HeapElement[T]
	Setter
}

type heapPrivate[T Lessor[T]] []T

func (h heapPrivate[T]) Len() int {
	return len(h)
}

func (h heapPrivate[T]) Less(i, j int) (ok bool) {
	ok = h[i].Less(h[j])
	return
}

func (h heapPrivate[T]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *heapPrivate[T]) Push(x any) {
	*h = append(*h, x.(T))
}

func (h *heapPrivate[T]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (a heapPrivate[T]) Copy() (b heapPrivate[T]) {
	b = heapPrivate[T](make([]T, 0, a.Len()))
	copy(b, a)
	return
}

func (a heapPrivate[T]) Sorted() (b heapPrivate[T]) {
	b = a.Copy()
	sort.Sort(b)
	return
}

func MakeHeap[T HeapElement[T]]() Heap[T] {
	return Heap[T]{
		l: &sync.Mutex{},
		h: heapPrivate[T](make([]T, 0)),
	}
}

func MakeHeapFromSlice[T HeapElement[T]](s heapPrivate[T]) Heap[T] {
	sort.Sort(s)

	return Heap[T]{
		l: &sync.Mutex{},
		h: s,
	}
}

type Heap[T HeapElement[T]] struct {
	l *sync.Mutex
	h heapPrivate[T]
	s int
}

func (h *Heap[T]) PeekPtr() (sk *T, ok bool) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.h.Len() > 0 {
		sk = &h.h[0]
		ok = true
	}

	return
}

func (h *Heap[T]) Peek() (sk T, ok bool) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.h.Len() > 0 {
		sk = h.h[0]
		ok = true
	}

	return
}

func (h *Heap[T]) Add(sk T) (err error) {
	h.Push(sk)
	return
}

func (h *Heap[T]) Push(sk T) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.s > 0 {
		panic(errors.Errorf("attempting to push to a heap that has saved elements"))
	}

	heap.Push(&h.h, sk)
}

func (h *Heap[T]) PopAndSave() (sk T, ok bool) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.h.Len() > 0 {
		sk = heap.Pop(&h.h).(T)
		ok = true
		h.s += 1
		faked := h.h[:h.h.Len()+h.s]
		faked[h.h.Len()] = sk
	}

	return
}

func (h *Heap[T]) Restore() {
	h.l.Lock()
	defer h.l.Unlock()

	h.h = h.h[:h.s]
	h.s = 0

	ReverseSortable(&h.h)

	return
}

func (h *Heap[T]) Pop() (sk T, ok bool) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.h.Len() > 0 {
		sk = heap.Pop(&h.h).(T)
		ok = true
	}

	return
}

func (h Heap[T]) Len() int {
	h.l.Lock()
	defer h.l.Unlock()

	return h.h.Len()
}

func (a Heap[T]) Equals(b Heap[T]) bool {
	a.l.Lock()
	defer a.l.Unlock()

	if a.h.Len() != b.h.Len() {
		return false
	}

	for i, av := range a.h {
		if b.h[i].Equals(&av) {
			return false
		}
	}

	return true
}

func (a Heap[T]) Copy() Heap[T] {
	a.l.Lock()
	defer a.l.Unlock()

	return Heap[T]{
		l: &sync.Mutex{},
		h: a.h.Copy(),
	}
}

func (a Heap[T]) EachPtr(f WriterFunc[*T]) (err error) {
	a.l.Lock()
	defer a.l.Unlock()

	for _, s := range a.h {
		if err = f(&s); err != nil {
			if IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a Heap[T]) Each(f WriterFunc[T]) (err error) {
	a.l.Lock()
	defer a.l.Unlock()

	for _, s := range a.h {
		if err = f(s); err != nil {
			if IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a Heap[T]) Sorted() (b heapPrivate[T]) {
	a.l.Lock()
	defer a.l.Unlock()

	b = a.h.Sorted()
	return
}

func (a *Heap[T]) Reset(b *Heap[T]) {
	a.l = &sync.Mutex{}

	if b == nil {
		a.h = heapPrivate[T](make([]T, 0))
	} else {
		a.h = heapPrivate[T](make([]T, b.Len()))

		for i, bv := range b.h {
			a.h[i] = bv
		}
	}
}
