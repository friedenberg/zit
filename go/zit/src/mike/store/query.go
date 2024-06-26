package store

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) Query(
	ms sku.QueryGroup,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.query(sku.QueryGroupWithKasten{QueryGroup: ms}, f, false)
}

func (s *Store) QueryOld(
	ms *query.Group,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.queryOld(query.GroupWithKasten{Group: ms}, f, false)
}

func (s *Store) QueryWithKasten(
	ms query.GroupWithKasten,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.queryOld(ms, f, true)
}

func (s *Store) queryOld(
	qg query.GroupWithKasten,
	f schnittstellen.FuncIter[*sku.Transacted],
	includeCwd bool,
) (err error) {
	if err = s.query(
		sku.QueryGroupWithKasten{
			QueryGroup: qg.Group,
			Kasten:     qg.Kasten,
		},
		f,
		includeCwd,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) query(
	qg sku.QueryGroupWithKasten,
	f schnittstellen.FuncIter[*sku.Transacted],
	includeCwd bool,
) (err error) {
	if qg.QueryGroup == nil {
		if qg.QueryGroup, err = s.queryBuilder.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var f1 schnittstellen.FuncIter[*sku.Transacted]

	// TODO improve performance by only reading Cwd zettels rather than scanning
	// everything
	if qg.GetSigil() == kennung.SigilCwd && includeCwd {
		f1 = func(z *sku.Transacted) (err error) {
			g := gattung.Must(z.GetGattung())
			m, ok := qg.Get(g)

			if !ok {
				return
			}

			if err = s.UpdateTransactedWithExternal(qg.Kasten, z); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !m.ContainsSku(z) {
				return
			}

			if err = f(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	} else {
		f1 = func(z *sku.Transacted) (err error) {
			g := gattung.Must(z.GetGattung())
			m, ok := qg.Get(g)

			if !ok {
				return
			}

			if includeCwd && m.GetSigil().IncludesCwd() {
				if err = s.UpdateTransactedWithExternal(qg.Kasten, z); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			if !m.ContainsSku(z) {
				return
			}

			if err = f(z); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	if err = s.GetVerzeichnisse().ReadQuery(
		qg.QueryGroup,
		f1,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// for _, vs := range s.virtualStores {
	// 	// TODO only query story if query group contains store type
	// 	if err = vs.Query(qg.Group, f1); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}
	// }

	return
}

func (s *Store) QueryCheckedOut(
	qg sku.ExternalQueryWithKasten,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	kid := qg.Kasten.GetKastenString()
	es, ok := s.externalStores[kid]

	if !ok {
		err = errors.Errorf("no kasten with id %q", kid)
		return
	}

	if err = es.QueryCheckedOut(
		qg.ExternalQuery,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryAllMatchingAkten(
	qg *query.Group,
	akten fd.Set,
	f func(*fd.FD, *sku.Transacted) error,
) (err error) {
	fds := fd.MakeMutableSetSha()

	var pa string

	if pa, err = s.GetStandort().DirObjektenGattung(
		s.GetKonfig().GetStoreVersion(),
		gattung.Akte,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = akten.Each(
		iter.MakeChain(
			func(fd *fd.FD) (err error) {
				if fd.GetShaLike().IsNull() {
					return iter.MakeErrStopIteration()
				}

				p := id.Path(fd.GetShaLike(), pa)

				if !files.Exists(p) {
					return iter.MakeErrStopIteration()
				}

				return
			},
			// TODO-P2 handle files with the same sha
			fds.Add,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	observed := fd.MakeMutableSet()
	var l sync.Mutex

	if err = s.GetVerzeichnisse().ReadQuery(
		qg,
		func(z *sku.Transacted) (err error) {
			fd, ok := fds.Get(z.GetAkteSha().String())

			if !ok {
				return
			}

			if err = f(fd, z); err != nil {
				err = errors.Wrap(err)
				return
			}

			l.Lock()
			defer l.Unlock()

			return observed.Add(fd)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fds.Each(
		func(fd *fd.FD) (err error) {
			if observed.Contains(fd) {
				return
			}

			if err = f(fd, nil); err != nil {
				err = errors.Wrap(err)
				return
			}

			return observed.Add(fd)
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
