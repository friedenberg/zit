package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) readExternalAndMergeIfNecessary(
	left, parent *sku.Transacted,
	options sku.CommitOptions,
) (err error) {
	if parent == nil {
		return
	}

	var col sku.CheckedOutLike

	if col, err = s.ReadCheckedOutFromTransacted(
		options.RepoId,
		parent,
	); err != nil {
		err = nil
		return
	}

	defer s.PutCheckedOutLike(col)

	right := col.GetSkuExternalLike().GetSku()

	parentEqualsExternal := right.Metadata.EqualsSansTai(&col.GetSku().Metadata)

	if parentEqualsExternal {
		op := checkout_options.Options{
			Force: true,
		}

		sku.TransactedResetter.ResetWith(right, left)

		if err = s.UpdateCheckoutFromCheckedOut(op, col); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	tm := sku.Conflicted{
		CheckedOutLike: col,
		Left:           left,
		Middle:         parent,
		Right:          right,
	}

	if err = s.Merge(tm); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Merge(
	tm sku.Conflicted,
) (err error) {
	switch tm.CheckedOutLike.GetRepoId().GetRepoIdString() {
	case "browser":
		err = todo.Implement()

	default:
		if err = s.cwdFiles.Merge(tm); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) RunMergeTool(
	tm sku.Conflicted,
) (err error) {
	tool := s.GetKonfig().ToolOptions.Merge

	switch tm.GetRepoId().GetRepoIdString() {
	case "browser":
		err = todo.Implement()

	default:
		var co sku.CheckedOutLike

		if co, err = s.cwdFiles.RunMergeTool(tool, tm); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer s.PutCheckedOutLike(co)

		if _, err = s.CreateOrUpdateCheckedOut(co, false); err != nil {
			err = errors.Wrap(err)
			return
		}

	}

	return
}
