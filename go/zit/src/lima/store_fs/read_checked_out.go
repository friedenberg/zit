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
	item *sku.FSItem,
) (co *sku.CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	// at a bare minimum, the internal object ID must always be set as there are
	// hard assumptions about internal being valid throughout the reading cycle
	co.Internal.ObjectId.ResetWith(&item.ExternalObjectId)

	if err = s.externalStoreSupplies.FuncReadOneInto(
		item.ExternalObjectId.String(),
		co.GetSku(),
	); err != nil {
		if collections.IsErrNotFound(err) || genres.IsErrUnsupportedGenre(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.HydrateExternalFromItem(
		sku.CommitOptions{
			Mode: object_mode.ModeUpdateTai,
		},
		item,
		co.GetSku(),
		&co.External,
	); err != nil {
		if errors.Is(err, sku.ErrExternalHasConflictMarker) {
			co.State = checked_out_state.Conflicted

			if err = co.External.ObjectId.SetWithIdLike(
				&co.Internal.ObjectId,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", item.Debug())
			return
		}
	}

	sku.DetermineState(co, false)

	if !item.Conflict.IsEmpty() {
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
	if co.GetSku() != sk {
		sku.Resetter.ResetWith(co.GetSku(), sk)
	}

	ok := false

	var kfp *sku.FSItem

	if kfp, ok = s.Get(&sk.ObjectId); !ok {
		err = collections.MakeErrNotFound(sk.GetObjectId())
		return
	}

	if err = s.HydrateExternalFromItem(
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
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", kfp)
		}

		return
	}

	sku.DetermineState(co, false)

	return
}
