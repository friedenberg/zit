package collections

import (
	"container/heap"
	"sort"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type HeapElement[T any] interface {
	schnittstellen.Equatable[T]
	Lessor[T]
}

type HeapElementPtr[T any] interface {
	HeapElement[T]
	schnittstellen.Resetable[T]
	schnittstellen.Ptr[T]
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
	b = heapPrivate[T](make([]T, a.Len()))
	copy(b, a)
	return
}

func (a heapPrivate[T]) Sorted() (b heapPrivate[T]) {
	b = a.Copy()
	sort.Sort(b)
	return
}

func MakeHeap[T HeapElement[T], T1 HeapElementPtr[T]]() Heap[T, T1] {
	return Heap[T, T1]{
		p: MakeFakePool[T, T1](),
		l: &sync.Mutex{},
		h: heapPrivate[T](make([]T, 0)),
	}
}

func MakeHeapFromSlice[T HeapElement[T], T1 HeapElementPtr[T]](
	s heapPrivate[T],
) Heap[T, T1] {
	sort.Sort(s)

	return Heap[T, T1]{
		p: MakeFakePool[T, T1](),
		l: &sync.Mutex{},
		h: s,
	}
}

type Heap[T HeapElement[T], T1 HeapElementPtr[T]] struct {
	p Pool2Like[T, T1]
	l *sync.Mutex
	h heapPrivate[T]
	s int
}

func (h *Heap[T, T1]) SetPool(p Pool2Like[T, T1]) {
	if p == nil {
		p = MakeFakePool[T, T1]()
	}

	h.p = p
}

func (h *Heap[T, T1]) Peek() (sk T1, ok bool) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.h.Len() > 0 {
		sk = h.p.Get()
		sk.ResetWith(h.h[0])
		ok = true
	}

	return
}

func (h *Heap[T, T1]) Add(sk T) (err error) {
	h.Push(sk)
	return
}

func (h *Heap[T, T1]) Push(sk T) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.s > 0 {
		panic(errors.Errorf("attempting to push to a heap that has saved elements"))
	}

	heap.Push(&h.h, sk)
}

func (h *Heap[T, T1]) PopAndSave() (sk T1, ok bool) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.h.Len() > 0 {
		sk = h.p.Get()
		sk.ResetWith(heap.Pop(&h.h).(T))
		ok = true
		h.s += 1
		faked := h.h[:h.h.Len()+h.s]
		faked[h.h.Len()] = *sk
	}

	return
}

func (h *Heap[T, T1]) Restore() {
	h.l.Lock()
	defer h.l.Unlock()

	h.h = h.h[:h.s]
	h.s = 0

	ReverseSortable(&h.h)

	return
}

func (h *Heap[T, T1]) Pop() (sk T1, ok bool) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.h.Len() > 0 {
		sk = h.p.Get()
		sk.ResetWith(heap.Pop(&h.h).(T))
		ok = true
	}

	return
}

func (h Heap[T, T1]) Len() int {
	h.l.Lock()
	defer h.l.Unlock()

	return h.h.Len()
}

func (a Heap[T, T1]) Equals(b Heap[T, T1]) bool {
	a.l.Lock()
	defer a.l.Unlock()

	if a.h.Len() != b.h.Len() {
		return false
	}

	for i, av := range a.h {
		if b.h[i].Equals(av) {
			return false
		}
	}

	return true
}

func (a Heap[T, T1]) Copy() Heap[T, T1] {
	a.l.Lock()
	defer a.l.Unlock()

	return Heap[T, T1]{
		l: &sync.Mutex{},
		h: a.h.Copy(),
	}
}

func (a Heap[T, T1]) EachPtr(f WriterFunc[T1]) (err error) {
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

func (a Heap[T, T1]) Each(f WriterFunc[T]) (err error) {
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

// func (a *Heap[T]) MergeStream(
//   read func() (T, error),
//   write WriterFunc[T],
// ) (err error) {
// 	defer func() {
// 		a.Restore()
// 	}()

// 	for {
// 		var e T

//     if e, err = read(); err != nil {
// 			if errors.IsEOF(err) {
// 				err = nil
// 				break
// 			} else {
// 				err = errors.Wrap(err)
// 				return
// 			}
//     }

// 	LOOP:
// 		for {
// 			peeked, ok := a.PeekPtr()

// 			switch {
// 			case !ok:
// 				break LOOP

// 			case peeked.Equals(e):
// 				a.Pop()
// 				continue

// 			case !peeked.Less(e):
// 				break LOOP

// 			default:
// 			}

// 			popped, _ := a.PopAndSave()

// 			if err = write(&popped); err != nil {
// 				if collections.IsStopIteration(err) {
// 					err = nil
// 				} else {
// 					err = errors.Wrap(err)
// 				}

// 				return
// 			}
// 		}

// 		if err = w(tz); err != nil {
// 			if collections.IsStopIteration(err) {
// 				err = nil
// 			} else {
// 				err = errors.Wrap(err)
// 			}

// 			return
// 		}
// 	}

// 	var last *zettel.Transacted

// 	for {
// 		popped, ok := zp.added.PopAndSave()

// 		if !ok {
// 			break
// 		}

// 		if last == nil {
// 			l := popped
// 			last = &l
// 		} else if popped.GetSku2().Less(last.GetSku2()) {
// 			err = errors.Errorf(
// 				"last time is greater than current! last: %s, current: %s, page: %d, less: %v, sku less: %v, sku2 less: %v",
// 				last.GetSku2(),
// 				popped.GetSku2(),
// 				zp.pageId.index,
// 				popped.Less(*last),
// 				popped.GetSku().Less(last.GetSku()),
// 				popped.GetSku2().Less(last.GetSku2()),
// 			)
// 			return
// 		}

// 		errors.Log().Printf(
// 			"page: %d post: %s time sku: %s time sku2: %s",
// 			zp.pageId.index,
// 			popped.GetSku2(),
// 			popped.Sku.GetTime(),
// 			popped.GetSku2().GetTime(),
// 		)

// 		if err = w(&popped); err != nil {
// 			if collections.IsStopIteration(err) {
// 				err = nil
// 			} else {
// 				err = errors.Wrap(err)
// 			}

// 			return
// 		}
// 	}
// 	return
// }

func (a *Heap[T, T1]) Fix() {
	a.l.Lock()
	defer a.l.Unlock()

	heap.Init(&a.h)
}

func (a *Heap[T, T1]) Sorted() (b heapPrivate[T]) {
	a.l.Lock()
	defer a.l.Unlock()

	b = a.h.Sorted()
	return
}

func (a *Heap[T, T1]) Reset() {
	a.l = &sync.Mutex{}
	a.h = heapPrivate[T](make([]T, 0))
	a.SetPool(nil)
}

func (a *Heap[T, T1]) ResetWith(b Heap[T, T1]) {
	a.l = &sync.Mutex{}

	a.h = heapPrivate[T](make([]T, b.Len()))

	for i, bv := range b.h {
		a.h[i] = bv
	}

	a.SetPool(b.p)
}
