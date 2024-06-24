package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) ReadOneExternalFS(
	sk2 *sku.Transacted,
) (co *sku.CheckedOutFS, err error) {
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

	var e2 *sku.ExternalFS

	if e2, err = s.ReadOneExternal(
		ObjekteOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		e,
		sk2,
	); err != nil {
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
	f schnittstellen.FuncIter[*sku.CheckedOutFS],
) (err error) {
	o := ObjekteOptions{
		Mode: objekte_mode.ModeRealize,
	}

	if err = s.cwdFiles.All(
		s.MakeHydrateExternalMaybe(qg, f, o),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) MakeHydrateExternalMaybe(
	qg *query.Group,
	f schnittstellen.FuncIter[*sku.CheckedOutFS],
	o ObjekteOptions,
) schnittstellen.FuncIter[*sku.ExternalMaybe] {
	return func(em *sku.ExternalMaybe) (err error) {
		if err = s.HydrateExternalMaybe(o, qg, em, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) HydrateExternalMaybe(
	o ObjekteOptions,
	qg *query.Group,
	em *sku.ExternalMaybe,
	f schnittstellen.FuncIter[*sku.CheckedOutFS],
) (err error) {
	var co *sku.CheckedOutFS

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

func (s *Store) ReadFilesUnsure(
	qg *query.Group,
	f schnittstellen.FuncIter[*sku.CheckedOutFS],
) (err error) {
	o := ObjekteOptions{
		Mode: objekte_mode.ModeRealize,
	}

	if err = s.cwdFiles.AllUnsure(
		s.MakeHydrateExternalMaybe(qg, f, o),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
