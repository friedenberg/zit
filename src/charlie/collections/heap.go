package collections

import (
	"container/heap"
	"sort"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

type HeapElement[T any] interface{}

type HeapElementPtr[T any] interface {
	schnittstellen.Resetable[T]
	schnittstellen.Ptr[T]
}

type heapPrivate[T HeapElement[T]] struct {
	less     func(T, T) bool
	elements []T
}

func (h heapPrivate[T]) Len() int {
	return len(h.elements)
}

func (h heapPrivate[T]) Less(i, j int) (ok bool) {
	ok = h.less(h.elements[i], h.elements[j])
	return
}

func (h heapPrivate[T]) Swap(i, j int) {
	h.elements[i], h.elements[j] = h.elements[j], h.elements[i]
}

func (h *heapPrivate[T]) Push(x any) {
	h.elements = append(h.elements, x.(T))
}

func (h *heapPrivate[T]) Pop() any {
	old := h.elements
	n := len(old)
	x := old[n-1]
	h.elements = old[0 : n-1]
	return x
}

func (a heapPrivate[T]) Copy() (b heapPrivate[T]) {
	b.less = a.less
	b.elements = make([]T, a.Len())
	copy(b.elements, a.elements)
	return
}

func (a heapPrivate[T]) Sorted() (b heapPrivate[T]) {
	b = a.Copy()
	sort.Sort(b)
	return
}

func MakeHeap[T HeapElement[T], T1 HeapElementPtr[T]](
	less func(T, T) bool,
	equals func(T, T) bool,
) Heap[T, T1] {
	return Heap[T, T1]{
		pool: MakeFakePool[T, T1](),
		lock: &sync.Mutex{},
		heap: heapPrivate[T]{
			less:     less,
			elements: make([]T, 0),
		},
		equals: equals,
	}
}

func MakeHeapFromSlice[T HeapElement[T], T1 HeapElementPtr[T]](
	s []T,
	less func(T, T) bool,
	equals func(T, T) bool,
) Heap[T, T1] {
	heap := heapPrivate[T]{
		elements: s,
		less:     less,
	}

	sort.Sort(heap)

	return Heap[T, T1]{
		pool:   MakeFakePool[T, T1](),
		lock:   &sync.Mutex{},
		equals: equals,
		heap:   heap,
	}
}

type Heap[T HeapElement[T], T1 HeapElementPtr[T]] struct {
	pool   schnittstellen.Pool[T, T1]
	lock   *sync.Mutex
	heap   heapPrivate[T]
	s      int
	equals func(T, T) bool
}

func (heap *Heap[T, T1]) SetPool(pool schnittstellen.Pool[T, T1]) {
	if pool == nil {
		pool = MakeFakePool[T, T1]()
	}

	heap.pool = pool
}

func (h *Heap[T, T1]) Peek() (sk T1, ok bool) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.heap.Len() > 0 {
		sk = h.pool.Get()
		sk.ResetWith(h.heap.elements[0])
		ok = true
	}

	return
}

func (h *Heap[T, T1]) Add(sk T) (err error) {
	h.Push(sk)
	return
}

func (h *Heap[T, T1]) Push(sk T) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.s > 0 {
		panic(
			errors.Errorf(
				"attempting to push to a heap that has saved elements",
			),
		)
	}

	heap.Push(&h.heap, sk)
}

func (h *Heap[T, T1]) PopAndSave() (sk T1, ok bool) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.heap.Len() > 0 {
		sk = h.pool.Get()
		sk.ResetWith(heap.Pop(&h.heap).(T))
		ok = true
		h.s += 1
		faked := h.heap.elements[:h.heap.Len()+h.s]
		faked[h.heap.Len()] = *sk
	}

	return
}

func (h *Heap[T, T1]) Restore() {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.heap.elements = h.heap.elements[:h.s]
	h.s = 0

	ReverseSortable(&h.heap)

	return
}

func (h *Heap[T, T1]) Pop() (sk T1, ok bool) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.heap.Len() > 0 {
		sk = h.pool.Get()
		sk.ResetWith(heap.Pop(&h.heap).(T))
		ok = true
	}

	return
}

func (h Heap[T, T1]) Len() int {
	h.lock.Lock()
	defer h.lock.Unlock()

	return h.heap.Len()
}

func (a Heap[T, T1]) Equals(b Heap[T, T1]) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.heap.Len() != b.heap.Len() {
		return false
	}

	for i, av := range a.heap.elements {
		if a.equals(b.heap.elements[i], av) {
			return false
		}
	}

	return true
}

func (a Heap[T, T1]) Copy() Heap[T, T1] {
	a.lock.Lock()
	defer a.lock.Unlock()

	return Heap[T, T1]{
		pool:   a.pool,
		lock:   &sync.Mutex{},
		equals: a.equals,
		heap:   a.heap.Copy(),
	}
}

func (a Heap[T, T1]) EachPtr(f schnittstellen.FuncIter[T1]) (err error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, s := range a.heap.elements {
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

func (a Heap[T, T1]) Each(f schnittstellen.FuncIter[T]) (err error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	for _, s := range a.heap.elements {
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

func (a *Heap[T, T1]) MergeStream(
	read func() (T1, error),
	write schnittstellen.FuncIter[T1],
) (err error) {
	defer func() {
		a.Restore()
	}()

	for {
		var e T1

		if e, err = read(); err != nil {
			if IsStopIteration(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

	LOOP:
		for {
			peeked, ok := a.Peek()

			switch {
			case !ok:
				break LOOP

			case a.equals(*peeked, *e):
				a.Pop()
				continue

			case !a.heap.less(*peeked, *e):
				break LOOP

			default:
			}

			popped, _ := a.PopAndSave()

			if err = write(popped); err != nil {
				if IsStopIteration(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}
		}

		if err = write(e); err != nil {
			if IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	var last T1

	for {
		popped, ok := a.PopAndSave()

		if !ok {
			break
		}

		if last == nil {
			last = popped
		} else if a.equals(*popped, *last) {
			continue
		} else if a.heap.less(*popped, *last) {
			err = errors.Errorf(
				"last is greater than current! last: %v, current: %v",
				last,
				popped,
			)
			return
		}

		if err = write(popped); err != nil {
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

func (a *Heap[T, T1]) Fix() {
	a.lock.Lock()
	defer a.lock.Unlock()

	heap.Init(&a.heap)
}

func (a *Heap[T, T1]) Sorted() (b1 []T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	b := a.heap.Sorted()
	b1 = b.elements

	return
}

func (a *Heap[T, T1]) Reset() {
	a.lock = &sync.Mutex{}
	a.heap.elements = make([]T, 0)
	a.SetPool(nil)
}

func (a *Heap[T, T1]) ResetWith(b Heap[T, T1]) {
	a.lock = &sync.Mutex{}

	a.heap.elements = make([]T, b.Len())
	a.heap.less = b.heap.less
	a.equals = b.equals

	for i, bv := range b.heap.elements {
		a.heap.elements[i] = bv
	}

	a.SetPool(b.pool)
}
