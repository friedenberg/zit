package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type Conflicted struct {
	CheckedOutLike
	Left, Middle, Right *Transacted
}

func (tm Conflicted) IsAllInlineTyp(itc kennung.InlineTypChecker) bool {
	if !itc.IsInlineTyp(tm.Left.GetTyp()) {
		return false
	}

	if !itc.IsInlineTyp(tm.Middle.GetTyp()) {
		return false
	}

	if !itc.IsInlineTyp(tm.Right.GetTyp()) {
		return false
	}

	return true
}

func (tm *Conflicted) MergeEtiketten() (err error) {
	left := tm.Left.GetEtiketten().CloneMutableSetPtrLike()
	middle := tm.Middle.GetEtiketten().CloneMutableSetPtrLike()
	right := tm.Right.GetEtiketten().CloneMutableSetPtrLike()

	same := kennung.MakeEtikettMutableSet()
	deleted := kennung.MakeEtikettMutableSet()

	removeFromAllButAddTo := func(
		e *kennung.Etikett,
		toAdd kennung.EtikettMutableSet,
	) (err error) {
		if err = toAdd.AddPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = left.DelPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = middle.DelPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = right.DelPtr(e); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = middle.EachPtr(
		func(e *kennung.Etikett) (err error) {
			if left.ContainsKey(left.KeyPtr(e)) && right.ContainsKey(right.KeyPtr(e)) {
				return removeFromAllButAddTo(e, same)
			} else if left.ContainsKey(left.KeyPtr(e)) || right.ContainsKey(right.KeyPtr(e)) {
				return removeFromAllButAddTo(e, deleted)
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = left.EachPtr(same.AddPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = right.EachPtr(same.AddPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	ets := same.CloneSetPtrLike()

	tm.Left.GetMetadatei().SetEtiketten(ets)
	tm.Middle.GetMetadatei().SetEtiketten(ets)
	tm.Right.GetMetadatei().SetEtiketten(ets)

	return
}

func (tm *Conflicted) ReadConflictMarker(
	s Scanner,
) (err error) {
	i := 0

	for s.Scan() {
		sk := s.GetTransacted()

		switch i {
		case 0:
			tm.Left = sk

		case 1:
			tm.Middle = sk

		case 2:
			tm.Right = sk

		default:
			err = errors.Errorf("too many skus in conflict file")
			return
		}

		i++
	}

	if err = s.Error(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (tm Conflicted) WriteConflictMarker(
	p ManyPrinter,
) (err error) {
	if _, err = p.PrintMany(tm.Left, tm.Middle, tm.Right); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
