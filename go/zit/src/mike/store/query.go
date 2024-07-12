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
	qg sku.QueryGroup,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if qg == nil {
		if qg, err = s.queryBuilder.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.GetVerzeichnisse().ReadQuery(
		qg,
		qg.MakeEmitSkuSigilSchwanzen(
			f,
			kennung.Kasten{},
			s.UpdateTransactedWithExternal,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryWithKasten(
	qg sku.ExternalQuery,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if qg.QueryGroup == nil {
		if qg.QueryGroup, err = s.queryBuilder.BuildQueryGroup(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	wg := iter.MakeErrorWaitGroupParallel()

	wg.Do(func() (err error) {
		if err = s.GetVerzeichnisse().ReadQuery(
			qg.QueryGroup,
			qg.QueryGroup.MakeEmitSkuIfNecessary(
				f,
				qg.Kasten,
				s.UpdateTransactedWithExternal,
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	})

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryCheckedOut(
	qg sku.ExternalQuery,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	kid := qg.Kasten.GetKastenString()
	es, ok := s.externalStores[kid]

	if !ok {
		err = errors.Errorf("no kasten with id %q", kid)
		return
	}

	if err = es.QueryCheckedOut(
		qg,
		f,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO [cot/gl !task project-2021-zit-kasten today zz-inbox] move unsure akten and untracked into kasten interface and store_fs
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
