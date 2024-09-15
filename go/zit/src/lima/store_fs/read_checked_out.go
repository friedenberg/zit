package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadCheckedOutFromObjectIdFDPair(
	o sku.CommitOptions,
	em *Item,
) (co *CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	if err = s.externalStoreInfo.FuncReadOneInto(
		em.ObjectId.String(),
		&co.Internal,
	); err != nil {
		if collections.IsErrNotFound(err) || genres.IsErrUnsupportedGenre(err) {
			// TODO mark status as new
			err = nil
			co.Internal.ObjectId.ResetWith(&em.ObjectId)
			co.State = checked_out_state.Untracked
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.ReadIntoCheckedOutFromTransactedAndFDSet(
		&co.Internal,
		em,
		co,
	); err != nil {
		if collections.IsErrNotFound(err) {
			// TODO mark status as new
			err = nil
			co.Internal.ObjectId.ResetWith(&em.ObjectId)
			co.State = checked_out_state.Untracked
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if !em.Conflict.IsEmpty() {
		co.State = checked_out_state.Conflicted
	}

	return
}

func (s *Store) ReadCheckedOutFromTransacted(
	sk2 *sku.Transacted,
) (co *CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	if err = s.ReadIntoCheckedOutFromTransacted(sk2, co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadIntoCheckedOutFromTransacted(
	sk *sku.Transacted,
	co *CheckedOut,
) (err error) {
	if &co.Internal != sk {
		sku.Resetter.ResetWith(&co.Internal, sk)
	}

	ok := false

	var kfp *Item

	if kfp, ok = s.Get(&sk.ObjectId); !ok {
		err = collections.MakeErrNotFound(sk.GetObjectId())
		return
	}

	if err = s.ReadIntoExternalFromObjectIdFDPair(
		sku.CommitOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		kfp,
		sk,
		&co.External,
	); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else if errors.Is(err, ErrExternalHasConflictMarker) {
			co.State = checked_out_state.Conflicted
			co.External.fds.ResetWith(kfp)

			if err = co.External.Transacted.ObjectId.SetWithIdLike(
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

func (s *Store) ReadIntoCheckedOutFromTransactedAndFDSet(
	sk *sku.Transacted,
	fds *Item,
	co *CheckedOut,
) (err error) {
	if &co.Internal != sk {
		sku.Resetter.ResetWith(&co.Internal, sk)
	}

	if err = s.ReadIntoExternalFromObjectIdFDPair(
		sku.CommitOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		fds,
		sk,
		&co.External,
	); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else if errors.Is(err, ErrExternalHasConflictMarker) {
			co.State = checked_out_state.Conflicted
			co.External.fds.ResetWith(fds)

			if err = co.External.Transacted.ObjectId.SetWithIdLike(
        &sk.ObjectId,
      ); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", fds)
		}

		return
	}

	sku.DetermineState(co, false)

	return
}
