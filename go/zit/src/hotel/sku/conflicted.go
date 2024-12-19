package sku

import (
	"iter"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type Conflicted struct {
	*CheckedOut
	Local, Base, Remote *Transacted
}

func (tm Conflicted) GetCollection() Collection {
	return tm
}

func (tm Conflicted) Len() int {
	return 3
}

func (tm Conflicted) Any() *Transacted {
	return tm.Local
}

func (c Conflicted) All() iter.Seq[*Transacted] {
	return func(yield func(*Transacted) bool) {
		if !yield(c.Local) {
			return
		}

		if c.Base != nil && !yield(c.Base) {
			return
		}

		if !yield(c.Remote) {
			return
		}
	}
}

func (tm Conflicted) IsAllInlineType(itc ids.InlineTypeChecker) bool {
	if !itc.IsInlineType(tm.Local.GetType()) {
		return false
	}

	if tm.Base != nil && !itc.IsInlineType(tm.Base.GetType()) {
		return false
	}

	if !itc.IsInlineType(tm.Remote.GetType()) {
		return false
	}

	return true
}

func (tm *Conflicted) MergeTags() (err error) {
	if tm.Base == nil {
		return
	}

	left := tm.Local.GetTags().CloneMutableSetPtrLike()
	middle := tm.Base.GetTags().CloneMutableSetPtrLike()
	right := tm.Remote.GetTags().CloneMutableSetPtrLike()

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

	tm.Local.GetMetadata().SetTags(ets)
	tm.Base.GetMetadata().SetTags(ets)
	tm.Remote.GetMetadata().SetTags(ets)

	return
}

func (tm *Conflicted) ReadConflictMarker(
	iter func(interfaces.FuncIter[*Transacted]) error,
) (err error) {
	i := 0

	if err = iter(
		func(sk *Transacted) (err error) {
			switch i {
			case 0:
				tm.Local = sk

			case 1:
				tm.Base = sk

			case 2:
				tm.Remote = sk

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

	// Conflicts can exist between objects without a base
	if i == 2 {
		tm.Base = tm.Remote
		tm.Remote = nil
	}

	return
}
