package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO remove Item from construction
func (s *Store) readCheckedOutFromItem(
	o sku.CommitOptions,
	i *sku.FSItem) (co *sku.CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	if err = s.externalStoreSupplies.FuncReadOneInto(
		i.ObjectId.String(),
		&co.Internal,
	); err != nil {
		if collections.IsErrNotFound(err) || genres.IsErrUnsupportedGenre(err) {
			err = nil
			co.Internal.ObjectId.ResetWith(&i.ObjectId)
			co.State = checked_out_state.Untracked
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.readIntoCheckedOutFromTransactedAndItem(
		&co.Internal,
		i,
		co,
	); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
			co.Internal.ObjectId.ResetWith(&i.ObjectId)
			co.State = checked_out_state.Untracked
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if !i.Conflict.IsEmpty() {
		co.State = checked_out_state.Conflicted
	}

	return
}

func (s *Store) ReadCheckedOutFromTransacted(
	sk2 *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	if err = s.readIntoCheckedOutFromTransacted(sk2, co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) readIntoCheckedOutFromTransacted(
	sk *sku.Transacted,
	co *sku.CheckedOut,
) (err error) {
	if &co.Internal != sk {
		sku.Resetter.ResetWith(&co.Internal, sk)
	}

	ok := false

	var kfp *sku.FSItem

	if kfp, ok = s.Get(&sk.ObjectId); !ok {
		err = collections.MakeErrNotFound(sk.GetObjectId())
		return
	}

	if err = s.readIntoExternalFromItem(
		sku.CommitOptions{
			Mode: object_mode.ModeUpdateTai,
		},
		kfp,
		sk,
		&co.External,
	); err != nil {
		if errors.IsNotExist(err) {
			err = quiter.MakeErrStopIteration()
		} else if errors.Is(err, sku.ErrExternalHasConflictMarker) {
			co.State = checked_out_state.Conflicted

			if err = co.External.ObjectId.SetWithIdLike(
				&sk.ObjectId,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", kfp)
		}

		return
	}

	sku.DetermineState(co, false)

	return
}

func (s *Store) readIntoCheckedOutFromTransactedAndItem(
	sk *sku.Transacted,
	i *sku.FSItem, co *sku.CheckedOut,
) (err error) {
	if &co.Internal != sk {
		sku.Resetter.ResetWith(&co.Internal, sk)
	}

	if err = s.readIntoExternalFromItem(
		sku.CommitOptions{
			Mode: object_mode.ModeUpdateTai,
		},
		i,
		sk,
		&co.External,
	); err != nil {
		if errors.IsNotExist(err) {
			err = errors.Wrapf(err, "Item: %s", i.Debug())
			return
			// err = iter.MakeErrStopIteration()
		} else if errors.Is(err, sku.ErrExternalHasConflictMarker) {
			co.State = checked_out_state.Conflicted

			if err = co.External.ObjectId.SetWithIdLike(
				&sk.ObjectId,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", i)
		}

		return
	}

	sku.DetermineState(co, false)

	return
}
