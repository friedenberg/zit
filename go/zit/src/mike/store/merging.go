package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) readExternalAndMergeIfNecessary(
	kinder, mutter *sku.Transacted,
	options sku.ObjekteOptions,
) (err error) {
	if mutter == nil {
		return
	}

	var col sku.CheckedOutLike

	if col, err = s.ReadCheckedOutFromTransacted(
		options.Kasten,
		mutter,
	); err != nil {
		err = nil
		return
	}

	defer s.PutCheckedOutLike(col)

	mutterEqualsExternal := sku.InternalAndExternalEqualsSansTai(col)

	if mutterEqualsExternal {
		op := checkout_options.Options{
			Force: true,
		}

		sku.TransactedResetter.ResetWith(col.GetSku(), kinder)

		if err = s.UpdateCheckoutFromCheckedOut(
			op,
			col,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		s.PutCheckedOutLike(col)

		return
	}

	transactedPtrCopy := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(transactedPtrCopy)

	if err = transactedPtrCopy.SetFromSkuLike(kinder); err != nil {
		err = errors.Wrap(err)
		return
	}

	tm := sku.Conflicted{
		CheckedOutLike: col,
		Left:           transactedPtrCopy,
		Middle:         col.GetSku(),
		Right:          col.GetSkuExternalLike().GetSku(),
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
	switch tm.CheckedOutLike.GetKasten().GetKastenString() {
	case "chrome":
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
