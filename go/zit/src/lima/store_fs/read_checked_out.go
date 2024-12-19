package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

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
		co.GetSkuExternal(),
	); err != nil {
		if errors.IsNotExist(err) {
			err = quiter.MakeErrStopIteration()
		} else if sku.IsErrMergeConflict(err) {
			co.SetState(checked_out_state.Conflicted)

			if err = co.GetSkuExternal().ObjectId.SetWithIdLike(
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

	return
}
