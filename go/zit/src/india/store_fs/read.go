package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadOneCheckedOut(
	o sku.ObjekteOptions,
	em *KennungFDPair,
) (co *CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	if err = s.storeFuncs.FuncReadOneInto(&em.Kennung, &co.Internal); err != nil {
		if collections.IsErrNotFound(err) {
			// TODO mark status as new
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.ReadOneKennungFDPairExternalInto(
		o,
		em,
		&co.Internal,
		&co.External,
	); err != nil {
		if errors.Is(err, ErrExternalHasConflictMarker) {
			err = nil
			co.State = checked_out_state.StateConflicted
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	sku.DetermineState(co, false)

	return
}

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

func (s *Store) ReadKennung(
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

	sku.DetermineState(co, false)

	return
}

func (s *Store) MakeHydrateCheckedOut(
	qg sku.Queryable,
	f schnittstellen.FuncIter[*CheckedOut],
	o sku.ObjekteOptions,
) schnittstellen.FuncIter[*KennungFDPair] {
	return func(em *KennungFDPair) (err error) {
		if err = s.HydrateCheckedOut(o, qg, em, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) HydrateCheckedOut(
	o sku.ObjekteOptions,
	qg sku.Queryable,
	em *KennungFDPair,
	f schnittstellen.FuncIter[*CheckedOut],
) (err error) {
	var co *CheckedOut

	if co, err = s.ReadOneCheckedOut(o, em); err != nil {
		err = errors.Wrapf(err, "%v", em)
		return
	}

	if !qg.ContainsSku(&co.External.Transacted) {
		return
	}

	if err = f(co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadQuery(
	qg sku.Queryable,
	f schnittstellen.FuncIter[*CheckedOut],
) (err error) {
	o := sku.ObjekteOptions{
		Mode: objekte_mode.ModeRealizeSansProto,
	}

	if err = s.All(
		s.MakeHydrateCheckedOut(qg, f, o),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadUnsure(
	qg sku.Queryable,
	f schnittstellen.FuncIter[*CheckedOut],
) (err error) {
	o := sku.ObjekteOptions{
		Mode: objekte_mode.ModeRealizeWithProto,
	}

	if err = s.AllUnsure(
		s.MakeHydrateCheckedOut(qg, f, o),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
