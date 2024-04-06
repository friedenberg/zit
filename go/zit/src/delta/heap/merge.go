package heap

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
)

func MergeStream[T Element, TPtr ElementPtr[T]](
	a *Heap[T, TPtr],
	read func() (TPtr, error),
	write schnittstellen.FuncIter[TPtr],
) (err error) {
	if err = MergeStreamPreferringHeap(
		a,
		read,
		write,
		a.h.equaler,
		a.h.Lessor,
		a.h.Resetter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MergeStreamPreferringHeap[T Element, TPtr ElementPtr[T]](
	h *Heap[T, TPtr],
	r func() (TPtr, error),
	w schnittstellen.FuncIter[TPtr],
	equaler schnittstellen.Equaler1[TPtr],
	l schnittstellen.Lessor3[TPtr],
	re schnittstellen.Resetter2[T, TPtr],
) (err error) {
	defer func() {
		h.Restore()
	}()

	oldWrite := w

	var last TPtr

	w = func(e TPtr) (err error) {
		if last == nil {
			var t T
			last = &t
		} else if equaler.Equals(e, last) {
			return
		} else if l.Less(e, last) {
			err = errors.Errorf(
				"last is greater than current! last:\n%v\ncurrent: %v",
				last,
				e,
			)

			return
		}

		re.ResetWith(last, e)

		return oldWrite(e)
	}

	for {
		var element TPtr

		if element, err = r(); err != nil {
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
			peeked, ok := h.Peek()

			switch {
			case !ok:
				break LOOP

			case equaler.Equals(peeked, element):
				element = peeked
				h.Pop()
				continue LOOP

			case l.Less(element, peeked):
				break LOOP

			default:
			}

			popped, _ := h.PopAndSave()

			if !equaler.Equals(peeked, popped) {
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

			if err = w(popped); err != nil {
				if iter.IsStopIteration(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}
		}

		if err = w(element); err != nil {
			if iter.IsStopIteration(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	for {
		popped, ok := h.PopAndSave()

		if !ok {
			break
		}

		if err = w(popped); err != nil {
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
