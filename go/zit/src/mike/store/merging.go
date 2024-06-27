package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO [radi/kof !task "add support for kasten in checkouts and external" project-2021-zit-features today zz-inbox]
func (s *Store) readExternalAndMergeIfNecessary(
	transactedPtr, mutter *sku.Transacted,
) (err error) {
	if mutter == nil {
		return
	}

	var col sku.CheckedOutLike

	if col, err = s.cwdFiles.ReadTransactedCheckedOut(mutter); err != nil {
		err = nil
		return
	}

	defer s.PutCheckedOutLike(col)

	mutterEqualsExternal := sku.InternalAndExternalEqualsSansTai(col)

	if mutterEqualsExternal {
		op := checkout_options.Options{
			Force: true,
		}

		if col, err = s.UpdateCheckout(
			col.GetKasten(),
			op,
			transactedPtr,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		s.PutCheckedOutLike(col)

		return
	}

	transactedPtrCopy := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(transactedPtrCopy)

	if err = transactedPtrCopy.SetFromSkuLike(transactedPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	tm := sku.Conflicted{
		CheckedOutLike: col,
		Left:           transactedPtrCopy,
		Middle:         col.GetSku(),
		Right:          col.GetSkuExternalLike().GetSku(),
	}

	if err = s.cwdFiles.Merge(tm); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) RunMergeTool(
	tm sku.Conflicted,
) (err error) {
	tool := s.GetKonfig().Cli().ToolOptions.Merge

	switch tm.GetKasten().GetKastenString() {
	case "chrome":
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
