package sku

import (
	"iter"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type Conflicted struct {
	CheckedOutLike
	Left, Middle, Right *Transacted
}

func (tm Conflicted) GetCollection() Collection {
	return tm
}

func (tm Conflicted) Len() int {
	return 3
}

func (tm Conflicted) Any() *Transacted {
	return tm.Left
}

func (tm Conflicted) All() iter.Seq[*Transacted] {
	return func(yield func(*Transacted) bool) {
		if !yield(tm.Left) {
			return
		}

		if !yield(tm.Middle) {
			return
		}

		if !yield(tm.Right) {
			return
		}
	}
}

func (tm Conflicted) IsAllInlineType(itc ids.InlineTypeChecker) bool {
	if !itc.IsInlineType(tm.Left.GetType()) {
		return false
	}

	if !itc.IsInlineType(tm.Middle.GetType()) {
		return false
	}

	if !itc.IsInlineType(tm.Right.GetType()) {
		return false
	}

	return true
}

func (tm *Conflicted) MergeTags() (err error) {
	left := tm.Left.GetTags().CloneMutableSetPtrLike()
	middle := tm.Middle.GetTags().CloneMutableSetPtrLike()
	right := tm.Right.GetTags().CloneMutableSetPtrLike()

	same := ids.MakeTagMutableSet()
	deleted := ids.MakeTagMutableSet()

	removeFromAllButAddTo := func(
		e *ids.Tag,
		toAdd ids.TagMutableSet,
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
		func(e *ids.Tag) (err error) {
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

	tm.Left.GetMetadata().SetTags(ets)
	tm.Middle.GetMetadata().SetTags(ets)
	tm.Right.GetMetadata().SetTags(ets)

	return
}

func (tm *Conflicted) ReadConflictMarker(
	iter func(interfaces.FuncIter[*Transacted]),
) (err error) {
	i := 0

	if iter(
		func(sk *Transacted) (err error) {
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

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
