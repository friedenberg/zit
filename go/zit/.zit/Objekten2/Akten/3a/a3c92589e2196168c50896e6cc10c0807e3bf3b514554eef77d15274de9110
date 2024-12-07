package store

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO-P2 add support for quiet reindexing
func (s *Store) Reindex() (err error) {
	if !s.GetDirectoryLayout().GetLockSmith().IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: "reindex",
		}

		return
	}

	if err = s.ResetIndexes(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetDirectoryLayout().ResetCache(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetStreamIndex().Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetInventoryListStore().ReadAllSkus(
		s.reindexOne,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CreateOrUpdate(
	in sku.ExternalLike,
	mode object_mode.Mode,
) (err error) {
	mode.Add(
		object_mode.ModeCommit,
		object_mode.ModeApplyProto,
	)

	if err = s.tryRealizeAndOrStore(
		in,
		sku.CommitOptions{
			Mode: mode,
		},
	); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return
}

func (s *Store) CreateOrUpdateBlobSha(
	k ids.IdLike,
	sh interfaces.Sha,
) (t *sku.Transacted, err error) {
	if !s.GetDirectoryLayout().GetLockSmith().IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				k.GetGenre(),
			),
		}

		return
	}

	t = sku.GetTransactedPool().Get()

	if err = t.ObjectId.SetWithIdLike(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadOneInto(k, t); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	t.SetBlobSha(sh)

	if err = s.tryRealizeAndOrStore(
		t,
		sku.CommitOptions{Mode: object_mode.ModeCommit},
	); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return
}

type RevertId struct {
	*ids.ObjectId
	ids.Tai
}

func (s *Store) RevertTo(
	ri RevertId,
) (err error) {
	if ri.Tai.IsEmpty() {
		return
	}

	if !s.GetDirectoryLayout().GetLockSmith().IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: "update many metadata",
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.GetStreamIndex().ReadOneObjectIdTai(
		ri.ObjectId,
		ri.Tai,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer sku.GetTransactedPool().Put(mutter)

	if err = s.tryRealizeAndOrStore(
		mutter,
		sku.CommitOptions{Mode: object_mode.ModeCommit},
	); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return
}

func (s *Store) CreateOrUpdateCheckedOut(
	col sku.SkuType,
	updateCheckout bool,
) (err error) {
	external := col.GetSkuExternal()
	internal := external.GetSku()

	if !s.GetDirectoryLayout().GetLockSmith().IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: fmt.Sprintf("create or update %s", internal.GetObjectId()),
		}

		return
	}

	if err = s.tryRealizeAndOrStore(
		external,
		sku.CommitOptions{Mode: object_mode.ModeCreate},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !updateCheckout {
		return
	}

	if err = s.UpdateCheckoutFromCheckedOut(
		checkout_options.OptionsWithoutMode{Force: true},
		col,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
