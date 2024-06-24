package store

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/id"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
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

			var e *sku.KennungFDPair

			e, ok = s.GetCwdFiles().Get(&z.Kennung)

			if !ok {
				return
			}

			var e2 *sku.ExternalFS

			if e2, err = s.ReadOneExternalFS(
				ObjekteOptions{
					Mode: objekte_mode.ModeUpdateTai,
				},
				e,
				z,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = z.SetFromSkuLike(&e2.Transacted); err != nil {
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
				var e *sku.KennungFDPair

				if e, ok = s.GetCwdFiles().Get(&z.Kennung); ok {
					var e2 *sku.ExternalFS

					if e2, err = s.ReadOneExternalFS(
						ObjekteOptions{
							Mode: objekte_mode.ModeUpdateTai,
						},
						e,
						z,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

					if err = z.SetFromSkuLike(&e2.Transacted); err != nil {
						err = errors.Wrap(err)
						return
					}
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

	if err = s.GetVerzeichnisse().ReadQuery(qg, f1); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, vs := range s.virtualStores {
		// TODO only query story if query group contains store type
		if err = vs.Query(qg, f1); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) ReadOne(
	k1 schnittstellen.StringerGattungGetter,
) (sk1 *sku.Transacted, err error) {
	sk1 = sku.GetTransactedPool().Get()

	if err = s.ReadOneInto(k1, sk1); err != nil {
		if collections.IsErrNotFound(err) {
			sku.GetTransactedPool().Put(sk1)
			sk1 = nil
		}

		err = errors.Wrap(err)
		return
	}

	return
}

// TODO [radi/kof !task "add support for kasten in checkouts and external" project-2021-zit-features today zz-inbox]
func (s *Store) ReadOneSigil(
	k1 schnittstellen.StringerGattungGetter,
	ka kennung.Kasten,
	si kennung.Sigil,
) (sk1 *sku.Transacted, err error) {
	sk1 = sku.GetTransactedPool().Get()

	if err = s.ReadOneInto(k1, sk1); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !si.IncludesExternal() {
		return
	}

	var k3 *kennung.Kennung3

	if k3, err = kennung.MakeKennung3(k1, ka); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ze sku.ExternalLike

	if ze, err = s.ReadOneExternal(
		ObjekteOptions{
			Mode: objekte_mode.ModeUpdateTai,
		},
		k3,
		sk1,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ze != nil {
		sku.TransactedResetter.ResetWith(sk1, ze.GetSku())
	}

	return
}

func (s *Store) ReadAllMatchingAkten(
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
