package heap

import (
	"container/heap"
	"sort"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/pool"
)

type Element interface{}

type ElementPtr[T Element] interface {
	schnittstellen.Ptr[T]
}

type heapPrivate[T Element, TPtr ElementPtr[T]] struct {
	Lessor   schnittstellen.Lessor3[TPtr]
	Resetter schnittstellen.Resetter2[T, TPtr]
	Elements []TPtr
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

func (h *heapPrivate[T, TPtr]) Pop() any {
	old := h.Elements
	n := len(old)
	x := old[n-1]
	h.Elements = old[0 : n-1]
	return x
}

func (a heapPrivate[T, TPtr]) Copy() (b heapPrivate[T, TPtr]) {
	b = heapPrivate[T, TPtr]{
		Lessor:   a.Lessor,
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

func Make[T Element, TPtr ElementPtr[T]](
	equaler schnittstellen.Equaler1[TPtr],
	lessor schnittstellen.Lessor3[TPtr],
	resetter schnittstellen.Resetter2[T, TPtr],
) Heap[T, TPtr] {
	return Heap[T, TPtr]{
		p:       pool.MakeFakePool[T, TPtr](),
		l:       &sync.Mutex{},
		equaler: equaler,
		h: heapPrivate[T, TPtr]{
			Lessor:   lessor,
			Resetter: resetter,
			Elements: make([]TPtr, 0),
		},
	}
}

func MakeHeapFromSlice[T Element, TPtr ElementPtr[T]](
	equaler schnittstellen.Equaler1[TPtr],
	lessor schnittstellen.Lessor3[TPtr],
	resetter schnittstellen.Resetter2[T, TPtr],
	s []TPtr,
) Heap[T, TPtr] {
	h := heapPrivate[T, TPtr]{
		Lessor:   lessor,
		Resetter: resetter,
		Elements: s,
	}

	sort.Sort(h)

	return Heap[T, TPtr]{
		p:       pool.MakeFakePool[T, TPtr](),
		l:       &sync.Mutex{},
		equaler: equaler,
		h:       h,
	}
}

type Heap[T Element, TPtr ElementPtr[T]] struct {
	p       schnittstellen.Pool[T, TPtr]
	l       *sync.Mutex
	h       heapPrivate[T, TPtr]
	equaler schnittstellen.Equaler1[TPtr]
	s       int
}

func (h *Heap[T, TPtr]) SetPool(p schnittstellen.Pool[T, TPtr]) {
	if p == nil {
		p = pool.MakeFakePool[T, TPtr]()
	}

	h.p = p
}

func (h *Heap[T, TPtr]) GetEqualer() schnittstellen.Equaler1[TPtr] {
	return h.equaler
}

func (h *Heap[T, TPtr]) GetLessor() schnittstellen.Lessor3[TPtr] {
	return h.h.Lessor
}

func (h *Heap[T, TPtr]) GetResetter() schnittstellen.Resetter2[T, TPtr] {
	return h.h.Resetter
}

func (h *Heap[T, TPtr]) Peek() (sk TPtr, ok bool) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.h.Len() > 0 {
		sk = h.p.Get()
		h.h.Resetter.ResetWithPtr(sk, h.h.Elements[0])
		ok = true
	}

	return
}

func (h *Heap[T, TPtr]) Add(sk TPtr) (err error) {
	h.Push(sk)
	return
}

func (h *Heap[T, TPtr]) Push(sk TPtr) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.s > 0 {
		panic(
			errors.Errorf(
				"attempting to push to a heap that has saved elements",
			),
		)
	}

	heap.Push(&h.h, sk)
}

func (h *Heap[T, TPtr]) PopAndSave() (sk TPtr, ok bool) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.h.Len() == 0 {
		return
	}

	sk = h.p.Get()
	e := heap.Pop(&h.h).(TPtr)
	h.h.Resetter.ResetWithPtr(sk, e)
	ok = true
	h.s += 1
	faked := h.h.Elements[:h.h.Len()+h.s]
	faked[h.h.Len()] = e

	return
}

func (h *Heap[T, TPtr]) Restore() {
	h.l.Lock()
	defer h.l.Unlock()

	h.h.Elements = h.h.Elements[:h.s]
	h.s = 0

	iter.ReverseSortable(&h.h)
}

func (h *Heap[T, TPtr]) PopError() (sk TPtr, err error) {
	ok := false
	sk, ok = h.Pop()

	if !ok {
		err = iter.MakeErrStopIteration()
	}

	return
}

func (h *Heap[T, TPtr]) Pop() (sk TPtr, ok bool) {
	h.l.Lock()
	defer h.l.Unlock()

	if h.h.Len() == 0 {
		return
	}

	sk = h.p.Get()
	h.h.Resetter.ResetWithPtr(sk, heap.Pop(&h.h).(TPtr))
	ok = true

	return
}

func (h Heap[T, TPtr]) Len() int {
	h.l.Lock()
	defer h.l.Unlock()

	return h.h.Len()
}

func (a Heap[T, TPtr]) Equals(b Heap[T, TPtr]) bool {
	a.l.Lock()
	defer a.l.Unlock()

	if a.h.Len() != b.h.Len() {
		return false
	}

	for i, av := range a.h.Elements {
		if b.equaler.Equals(b.h.Elements[i], av) {
			return false
		}
	}

	return true
}

func (a Heap[T, TPtr]) Copy() Heap[T, TPtr] {
	a.l.Lock()
	defer a.l.Unlock()

	return Heap[T, TPtr]{
		p:       a.p,
		equaler: a.equaler,
		l:       &sync.Mutex{},
		h:       a.h.Copy(),
	}
}

func (a Heap[T, TPtr]) EachPtr(f schnittstellen.FuncIter[TPtr]) (err error) {
	a.l.Lock()
	defer a.l.Unlock()

	for _, s := range a.h.Elements {
		if err = f(s); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a Heap[T, TPtr]) Each(f schnittstellen.FuncIter[T]) (err error) {
	a.l.Lock()
	defer a.l.Unlock()

	for _, s := range a.h.Elements {
		if err = f(*s); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a *Heap[T, TPtr]) MergeStream(
	read func() (TPtr, error),
	write schnittstellen.FuncIter[TPtr],
) (err error) {
	defer func() {
		a.Restore()
	}()

	for {
		var e TPtr

		if e, err = read(); err != nil {
			if iter.IsStopIteration(err) {
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

			case a.equaler.Equals(peeked, e):
				a.Pop()
				continue LOOP

			case !a.h.Lessor.Less(peeked, e):
				break LOOP

			default:
			}

			popped, _ := a.PopAndSave()

			if err = write(popped); err != nil {
				if iter.IsStopIteration(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}
		}

		if err = write(e); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	var last TPtr

	for {
		popped, ok := a.PopAndSave()

		if !ok {
			break
		}

		if last == nil {
			last = popped
		} else if a.equaler.Equals(popped, last) {
			continue
		} else if a.h.Lessor.Less(popped, last) {
			err = errors.Errorf(
				"last is greater than current! last: %v, current: %v",
				last,
				popped,
			)
			return
		}

		if err = write(popped); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	return
}

func (a *Heap[T, TPtr]) Fix() {
	a.l.Lock()
	defer a.l.Unlock()

	heap.Init(&a.h)
}

func (a *Heap[T, TPtr]) Sorted() (b []TPtr) {
	a.l.Lock()
	defer a.l.Unlock()

	b = a.h.Sorted().Elements

	return
}

func (a *Heap[T, TPtr]) Reset() {
	a.l = &sync.Mutex{}
	a.h.Elements = make([]TPtr, 0)
	a.SetPool(nil)
}

func (a *Heap[T, TPtr]) ResetWith(b Heap[T, TPtr]) {
	a.l = &sync.Mutex{}

	a.equaler = b.equaler

	a.h.Lessor = b.h.Lessor
	a.h.Resetter = b.h.Resetter
	a.h.Elements = make([]TPtr, b.Len())

	for i, bv := range b.h.Elements {
		a.h.Elements[i] = bv
	}

	a.SetPool(b.p)
}
