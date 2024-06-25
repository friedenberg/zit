package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

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
