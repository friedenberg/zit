package store

import (
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/files"
	"code.linenisgreat.com/zit/src/bravo/id"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/query"
)

func (s *Store) QueryWithoutCwd(
	ms *query.Group,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.query(ms, f, false)
}

func (s *Store) QueryWithCwd(
	ms *query.Group,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.query(ms, f, true)
}

func (s *Store) query(
	qg *query.Group,
	f schnittstellen.FuncIter[*sku.Transacted],
	includeCwd bool,
) (err error) {
	f1 := func(z *sku.Transacted) (err error) {
		g := gattung.Must(z.GetGattung())
		m, ok := qg.Get(g)

		if !ok {
			return
		}

		if includeCwd && m.GetSigil().IncludesCwd() {
			var e *sku.ExternalMaybe

			if e, ok = s.GetCwdFiles().Get(&z.Kennung); ok {
				var e2 *sku.External

				if e2, err = s.ReadOneExternal(e, z); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = z.SetFromSkuLike(&e2.Transacted); err != nil {
					err = errors.Wrap(err)
					return
				}
			}
		}

		if !m.ContainsMatchable(z) {
			return
		}

		return f(z)
	}

	if err = s.ReadQuery(qg, f1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOne(
	k1 schnittstellen.StringerGattungGetter,
) (sk1 *sku.Transacted, err error) {
	sk1 = sku.GetTransactedPool().Get()

	if err = s.ReadOneInto(k1, sk1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadQuery(
	qg *query.Group,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.GetVerzeichnisse().ReadQuery(qg, f)
}

func (s Store) ReadAllMatchingAkten(
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

	if err = s.ReadQuery(
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
