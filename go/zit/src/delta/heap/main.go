package heap

import (
	"container/heap"
	"sync"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/bravo/iter"
)

type Heap[T Element, TPtr ElementPtr[T]] struct {
	l sync.Mutex
	h heapPrivate[T, TPtr]
	s int
}

func (h *Heap[T, TPtr]) SetPool(v schnittstellen.Pool[T, TPtr]) {
	h.h.p = v
}

func (h *Heap[T, TPtr]) GetEqualer() schnittstellen.Equaler1[TPtr] {
	return h.h.equaler
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
		sk = h.h.GetPool().Get()
		h.h.Resetter.ResetWith(sk, h.h.Elements[0])
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

	// h.h.discardDupes()

	if h.h.Len() == 0 {
		return
	}

	sk = h.h.GetPool().Get()
	e := heap.Pop(&h.h).(TPtr)
	h.h.Resetter.ResetWith(sk, e)
	ok = true
	h.s += 1
	faked := h.h.Elements[:h.h.Len()+h.s]
	faked[h.h.Len()] = e
	h.h.saveLastPopped(sk)

	return
}

func (h *Heap[T, TPtr]) Restore() {
	h.l.Lock()
	defer h.l.Unlock()

	h.h.Elements = h.h.Elements[:h.s]
	h.s = 0
	h.h.GetPool().Put(h.h.lastPopped)
	h.h.lastPopped = nil

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

	// h.h.discardDupes()

	if h.h.Len() == 0 {
		return
	}

	sk = h.h.GetPool().Get()
	h.h.Resetter.ResetWith(sk, heap.Pop(&h.h).(TPtr))
	ok = true
	h.h.saveLastPopped(sk)

	return
}

func (h *Heap[T, TPtr]) Len() int {
	h.l.Lock()
	defer h.l.Unlock()

	return h.h.Len()
}

func (a *Heap[T, TPtr]) Equals(b *Heap[T, TPtr]) bool {
	a.l.Lock()
	defer a.l.Unlock()

	if a.h.Len() != b.h.Len() {
		return false
	}

	for i, av := range a.h.Elements {
		if b.h.equaler.Equals(b.h.Elements[i], av) {
			return false
		}
	}

	return true
}

func (a *Heap[T, TPtr]) Copy() Heap[T, TPtr] {
	a.l.Lock()
	defer a.l.Unlock()

	return Heap[T, TPtr]{
		h: a.h.Copy(),
	}
}

func (a *Heap[T, TPtr]) EachPtr(f schnittstellen.FuncIter[TPtr]) (err error) {
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

func (a *Heap[T, TPtr]) Each(f schnittstellen.FuncIter[T]) (err error) {
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

	oldWrite := write

	var last TPtr

	write = func(e TPtr) (err error) {
		if last == nil {
			var t T
			last = &t
		} else if a.h.equaler.Equals(e, last) {
			return
		} else if a.h.Lessor.Less(e, last) {
			err = errors.Errorf(
				"last is greater than current! last:\n%v\ncurrent: %v",
				last,
				e,
			)

			return
		}

		a.h.Resetter.ResetWith(last, e)

		return oldWrite(e)
	}

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

			case a.h.equaler.Equals(peeked, e):
				e = peeked
				a.Pop()
				continue LOOP

			case a.h.Lessor.Less(e, peeked):
				break LOOP

			default:
			}

			popped, _ := a.PopAndSave()

			if !a.h.equaler.Equals(peeked, popped) {
				err = errors.Errorf(
					"popped not equal to peeked: %s != %s",
					popped,
					peeked,
				)

				return
			}

			if popped == nil {
				break
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

		if err = write(e); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	for {
		popped, ok := a.PopAndSave()

		if !ok {
			break
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
	a.h.Elements = make([]TPtr, 0)
	a.h.GetPool().Put(a.h.lastPopped)
	a.h.p = nil
	a.h.lastPopped = nil
}

func (a *Heap[T, TPtr]) ResetWith(b *Heap[T, TPtr]) {
	a.h.equaler = b.h.equaler
	a.h.Lessor = b.h.Lessor
	a.h.Resetter = b.h.Resetter
	a.h.Elements = make([]TPtr, b.Len())

	for i, bv := range b.h.Elements {
		a.h.Elements[i] = bv
	}

	a.h.p = b.h.GetPool()
}
