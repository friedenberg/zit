package store

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/query"
)

func (s *Store) ReadOneExternalFS(
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

func (s *Store) ReadFiles(
	qg *query.Group,
	f schnittstellen.FuncIter[*sku.CheckedOut],
) (err error) {
	if err = s.cwdFiles.All(
		func(em *sku.ExternalMaybe) (err error) {
			var co *sku.CheckedOut

			if co, err = s.ReadOneCheckedOut(em); err != nil {
				err = errors.Wrap(err)
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
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
