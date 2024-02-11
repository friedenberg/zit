package store_util

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/matcher"
)

func (s *common) ReadOneExternalFS(
	sk2 *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	co = sku.GetCheckedOutPool().Get()

	if err = co.Internal.SetFromSkuLike(sk2); err != nil {
		err = errors.Wrap(err)
		return
	}

	ok := false

	var e *sku.ExternalMaybe

	if e, ok = s.cwdFiles.Get(&sk2.Kennung); !ok {
		err = iter.MakeErrStopIteration()
		return
	}

	var e2 *sku.External

	if e2, err = s.ReadOneExternal(e, sk2); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else if errors.Is(err, sku.ErrExternalHasConflictMarker) {
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

func (s *common) ReadFiles(
	fq matcher.FuncReaderTransactedLikePtr,
	f schnittstellen.FuncIter[*sku.CheckedOut],
) (err error) {
	if err = fq(
		iter.MakeChain(
			func(et *sku.Transacted) (err error) {
				var col *sku.CheckedOut

				et1 := sku.GetTransactedPool().Get()

				if err = et1.SetFromSkuLike(et); err != nil {
					err = errors.Wrap(err)
					return
				}

				if col, err = s.ReadOneExternalFS(et1); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = col.Internal.SetFromSkuLike(et); err != nil {
					err = errors.Wrap(err)
					return
				}

				if col.State == checked_out_state.StateUnknown {
					col.DetermineState(false)
				}

				if err = f(col); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.cwdFiles.EachCreatableMatchable(
		iter.MakeChain(
			func(il *sku.ExternalMaybe) (err error) {
				if err = s.GetAbbrStore().Exists(&il.Kennung); err == nil {
					err = iter.MakeErrStopIteration()
					return
				}

				err = nil

				tco := &sku.CheckedOut{}
				var tcoe *sku.External

				if tcoe, err = s.ReadOneExternal(
					il,
					nil,
				); err != nil {
					if errors.IsNotExist(err) {
						err = iter.MakeErrStopIteration()
					} else {
						err = errors.Wrapf(err, "%#v", il)
					}

					return
				}

				tco.Internal.Kennung.ResetWithKennung(&tcoe.Kennung)
				tco.External.SetFromSkuLike(tcoe)
				tco.State = checked_out_state.StateUntracked

				if err = f(tco); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
