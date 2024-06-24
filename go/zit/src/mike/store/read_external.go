package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/delta/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

// zips a CheckedOutFS from a known internal Sku with whatever Sku that may be
// checked out. If there is no checked, returns ErrStopIteration
//
// TODO [radi/kof !task "add support for kasten in checkouts and external" project-2021-zit-features today zz-inbox]
func (s *Store) CombineOneCheckedOutFS(
	sk2 *sku.Transacted,
) (co *store_fs.CheckedOut, err error) {
	co = store_fs.GetCheckedOutPool().Get()

	if err = co.Internal.SetFromSkuLike(sk2); err != nil {
		err = errors.Wrap(err)
		return
	}

	ok := false

	var e *store_fs.KennungFDPair

	if e, ok = s.cwdFiles.Get(&sk2.Kennung); !ok {
		err = iter.MakeErrStopIteration()
		return
	}

	var e2 *store_fs.External

	if e2, err = s.ReadOneExternalFS(
		ObjekteOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		e,
		sk2,
	); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else if errors.Is(err, store_fs.ErrExternalHasConflictMarker) {
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

// TODO [radi/kof !task "add support for kasten in checkouts and external" project-2021-zit-features today zz-inbox]
func (s *Store) ReadExternal(
	qg query.GroupWithKasten,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	switch qg.Kasten.String() {
	case "chrome":
		err = todo.Implement()

	default:
		if err = s.ReadExternalFS(
			qg.Group,
			func(cofs *store_fs.CheckedOut) (err error) {
				return f(cofs)
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) ReadExternalFS(
	qg *query.Group,
	f schnittstellen.FuncIter[*store_fs.CheckedOut],
) (err error) {
	o := ObjekteOptions{
		Mode: objekte_mode.ModeRealizeSansProto,
	}

	if err = s.cwdFiles.All(
		s.MakeHydrateCheckedOutFS(qg, f, o),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) MakeHydrateCheckedOutFS(
	qg *query.Group,
	f schnittstellen.FuncIter[*store_fs.CheckedOut],
	o ObjekteOptions,
) schnittstellen.FuncIter[*store_fs.KennungFDPair] {
	return func(em *store_fs.KennungFDPair) (err error) {
		if err = s.HydrateCheckedOutFS(o, qg, em, f); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func (s *Store) HydrateCheckedOutFS(
	o ObjekteOptions,
	qg *query.Group,
	em *store_fs.KennungFDPair,
	f schnittstellen.FuncIter[*store_fs.CheckedOut],
) (err error) {
	var co *store_fs.CheckedOut

	if co, err = s.ReadOneCheckedOutFS(o, em); err != nil {
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

// TODO [radi/kof !task "add support for kasten in checkouts and external" project-2021-zit-features today zz-inbox]
func (s *Store) ReadExternalFSUnsure(
	qg *query.Group,
	f schnittstellen.FuncIter[*store_fs.CheckedOut],
) (err error) {
	o := ObjekteOptions{
		Mode: objekte_mode.ModeRealizeWithProto,
	}

	if err = s.cwdFiles.AllUnsure(
		s.MakeHydrateCheckedOutFS(qg, f, o),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
