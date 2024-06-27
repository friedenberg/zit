package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadCheckedOutFromKennungFDPair(
	o sku.ObjekteOptions,
	em *KennungFDPair,
) (co *CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	if err = s.storeFuncs.FuncReadOneInto(&em.Kennung, &co.Internal); err != nil {
		if collections.IsErrNotFound(err) {
			// TODO mark status as new
			err = nil
			co.Internal.Kennung.ResetWith(&em.Kennung)
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.ReadIntoCheckedOutFromTransacted(&co.Internal, co); err != nil {
		err = errors.Wrap(err)
		return
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
		if err = co.Internal.SetFromSkuLike(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	ok := false

	var kfp *KennungFDPair

	if kfp, ok = s.Get(&sk.Kennung); !ok {
		err = iter.MakeErrStopIteration()
		return
	}

	if err = s.ReadIntoExternalFromKennungFDPair(
		sku.ObjekteOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		kfp,
		sk,
		&co.External,
	); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else if errors.Is(err, ErrExternalHasConflictMarker) {
			co.State = checked_out_state.StateConflicted
			co.External.FDs = kfp.FDs

			if err = co.External.Kennung.SetWithKennung(&sk.Kennung); err != nil {
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
