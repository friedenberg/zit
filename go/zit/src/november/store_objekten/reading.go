package store_objekten

import (
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/files"
	"code.linenisgreat.com/zit/src/bravo/id"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/fd"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/matcher_proto"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func (s *Store) QueryWithoutCwd(
	ms matcher_proto.QueryGroup,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.query(ms, f, false)
}

func (s *Store) QueryWithCwd(
	ms matcher_proto.QueryGroup,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.query(ms, f, true)
}

func (s *Store) query(
	ms matcher_proto.QueryGroup,
	f schnittstellen.FuncIter[*sku.Transacted],
	includeCwd bool,
) (err error) {
	gsWithoutHistory, gsWithHistory := matcher_proto.SplitGattungenByHistory(ms)

	wg := iter.MakeErrorWaitGroupParallel()

	f1 := func(z *sku.Transacted) (err error) {
		g := gattung.Must(z.GetGattung())
		m, ok := ms.Get(g)

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

	wg.Do(
		func() error {
			return s.ReadAllSchwanzen(gsWithoutHistory, f1)
		},
	)

	wg.Do(
		func() error {
			return s.ReadAll(gsWithHistory, f1)
		},
	)

	return wg.GetError()
}

func (s *Store) ReadOneInto(
	k1 schnittstellen.StringerGattungGetter,
	out *sku.Transacted,
) (err error) {
	var sk *sku.Transacted

	switch k1.GetGattung() {
	case gattung.Zettel:
		var h kennung.Hinweis

		if err = h.Set(k1.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if sk, err = s.ReadOneKennung(h); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Typ:
		var k kennung.Typ

		if err = k.Set(k1.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		sk = s.StoreUtil.GetKonfig().GetApproximatedTyp(k).ActualOrNil()

	case gattung.Etikett:
		var e kennung.Etikett

		if err = e.Set(k1.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		ok := false
		sk, ok = s.StoreUtil.GetKonfig().GetEtikett(e)

		if !ok {
			sk = nil
		}

	case gattung.Kasten:
		var k kennung.Kasten

		if err = k.Set(k.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		sk = s.StoreUtil.GetKonfig().GetKasten(k)

	case gattung.Konfig:
		sk = &s.StoreUtil.GetKonfig().Sku

		if sk.GetTai().IsEmpty() {
			sk = nil
		}

	default:
		err = errors.Errorf("unsupported gattung: %q -> %q", k1.GetGattung(), k1)
		return
	}

	if sk == nil {
		err = collections.MakeErrNotFound(k1)
		return
	}

	if err = out.SetFromSkuLike(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.virtualStore.HydrateOneChrome(sk); err != nil {
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

func (s *Store) MakeReadAllSchwanzen(
	gs ...gattung.Gattung,
) func(schnittstellen.FuncIter[*sku.Transacted]) error {
	return func(f schnittstellen.FuncIter[*sku.Transacted]) (err error) {
		return s.ReadAllSchwanzen(kennung.MakeGattung(gs...), f)
	}
}

func (s *Store) ReadAllSchwanzen(
	gs kennung.Gattung,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.GetVerzeichnisse().ReadSchwanzen(
		iter.MakeChain(
			func(sk *sku.Transacted) (err error) {
				if !gs.Contains(gattung.Must(sk.Kennung.GetGattung())) {
					err = iter.MakeErrStopIteration()
					return
				}

				return
			},
			s.virtualStore.HydrateOneChrome,
			f,
		),
	)
}

func (s *Store) ReadAll(
	gs kennung.Gattung,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.ReadAllGattungenFromVerzeichnisse(gs, f)
}

func (s Store) ReadAllMatchingAkten(
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

	if err = s.ReadAll(
		kennung.MakeGattung(gattung.TrueGattung()...),
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
