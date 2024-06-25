package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadOneKennungFDPairExternal(
	o sku.ObjekteOptions,
	em *KennungFDPair,
	t *sku.Transacted,
) (e *External, err error) {
	e = GetExternalPool().Get()

	if err = s.ReadOneKennungFDPairExternalInto(o, em, t, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneKennungFDPairExternalInto(
	o sku.ObjekteOptions,
	em *KennungFDPair,
	t *sku.Transacted,
	e *External,
) (err error) {
	o.Del(objekte_mode.ModeApplyProto)

	if err = s.ReadOneExternalInto(
		&o,
		em,
		t,
		e,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.storeFuncs.FuncCommit(
		&e.Transacted,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneKennungExternalFS(
	o sku.ObjekteOptions,
	k1 schnittstellen.StringerGattungGetter,
	sk *sku.Transacted,
) (el *External, err error) {
	e, ok := s.Get(k1)

	if !ok {
		return
	}

	if el, err = s.ReadOneKennungFDPairExternal(o, e, sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadTransactedCheckedOut(
	sk2 *sku.Transacted,
) (co *CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	if err = co.Internal.SetFromSkuLike(sk2); err != nil {
		err = errors.Wrap(err)
		return
	}

	ok := false

	var e *KennungFDPair

	if e, ok = s.Get(&sk2.Kennung); !ok {
		err = iter.MakeErrStopIteration()
		return
	}

	var e2 *External

	if e2, err = s.ReadOneKennungFDPairExternal(
		sku.ObjekteOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		e,
		sk2,
	); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else if errors.Is(err, ErrExternalHasConflictMarker) {
			co.State = checked_out_state.StateConflicted
			co.External.FDs = e.FDs

			if err = co.External.Kennung.SetWithKennung(&sk2.Kennung); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", e)
		}

		return
	}

	if err = co.External.SetFromSkuLike(e2); err != nil {
		err = errors.Wrap(err)
		return
	}

	co.DetermineState(false)

	return
}

func (s *Store) ReadOneKennung(
	o sku.ObjekteOptions,
	k1 schnittstellen.StringerGattungGetter,
	t *sku.Transacted,
) (e *External, err error) {
	k, ok := s.Get(k1)

	if !ok {
		return
	}

	if e, err = s.ReadOneKennungFDPairExternal(o, k, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
