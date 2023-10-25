package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/iter2"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

func (s *Store) onNewOrUpdated(
	t *sku.Transacted,
) (err error) {
	return s.onNewOrUpdatedCommit(t, true)
}

func (s *Store) onNewOrUpdatedCommit(
	t *sku.Transacted,
	commit bool,
) (err error) {
	if commit {
		s.StoreUtil.CommitUpdatedTransacted(t)
	}

	g := gattung.Must(t.Kennung.GetGattung())

	switch g {
	case gattung.Konfig:
		if err = s.StoreUtil.GetKonfigPtr().SetTransacted(t, s.GetAkten().GetKonfigV0()); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Kasten:
		if err = s.StoreUtil.GetKonfigPtr().AddKasten(t); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Typ:
		if err = s.StoreUtil.GetKonfigPtr().AddTyp(t); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Etikett:
		if err = s.StoreUtil.GetKonfigPtr().AddEtikett(t); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Zettel:
		if err = s.zettelStore.writeNamedZettelToIndex(t); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = gattung.MakeErrUnsupportedGattung(g)
		return
	}

	return
}

func (s *Store) onNew(
	t *sku.Transacted,
) (err error) {
	if err = s.onNewOrUpdated(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return s.LogWriter.New(t)
}

func (s *Store) onUpdated(
	t *sku.Transacted,
) (err error) {
	if err = s.onNewOrUpdated(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return s.LogWriter.Updated(t)
}

func (s *Store) onUnchanged(
	t *sku.Transacted,
) (err error) {
	return s.LogWriter.Unchanged(t)
}

func (s *Store) ReadOne(
	k1 schnittstellen.StringerGattungGetter,
) (sk1 *sku.Transacted, err error) {
	var sk *sku.Transacted

	switch k1.GetGattung() {
	case gattung.Zettel:
		var h kennung.Hinweis

		if err = h.Set(k1.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		if sk, err = s.zettelStore.verzeichnisseSchwanzen.ReadHinweisSchwanzen(h); err != nil {
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
		if sk, err = s.konfigStore.ReadOne(kennung.Konfig{}); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("unsupported gattung: %q -> %q", k1.GetGattung(), k1)
		return
	}

	if sk == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k1})
		return
	}

	sk1 = sku.GetTransactedPool().Get()

	if err = sk1.SetFromSkuLike(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadAllSchwanzen(
	gs gattungen.Set,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if err = iter2.Parallel(
		gs,
		func(g gattung.Gattung) (err error) {
			switch g {
			case gattung.Typ:
				if err = s.StoreUtil.GetKonfig().Typen.EachPtr(f); err != nil {
					err = errors.Wrap(err)
					return
				}

			case gattung.Etikett:
				if err = s.StoreUtil.GetKonfig().EachEtikett(f); err != nil {
					err = errors.Wrap(err)
					return
				}

			case gattung.Kasten:
				if err = s.GetKonfig().Kisten.EachPtr(f); err != nil {
					err = errors.Wrap(err)
					return
				}

			case gattung.Konfig:
				var k *sku.Transacted

				if k, err = s.konfigStore.ReadOne(&kennung.Konfig{}); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = f(k); err != nil {
					err = errors.Wrap(err)
					return
				}

			case gattung.Zettel:
				return s.zettelStore.verzeichnisseSchwanzen.ReadMany(f)

			default:
				err = todo.Implement()
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadAll(
	gs gattungen.Set,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.ReadAllGattungen(gs, f)
}
