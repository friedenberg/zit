package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/objekte_mode"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/charlie/hinweisen"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

func (s *Store) handleNewOrUpdated(
	t *sku.Transacted,
	updateType objekte_mode.Mode,
) (err error) {
	return iter.Chain(
		t,
		s.AddMatchable,
		func(t *sku.Transacted) error {
			return s.handleNewOrUpdatedCommit(t, updateType)
		},
	)
}

func (s *Store) handleNewOrUpdatedCommit(
	t *sku.Transacted,
	mode objekte_mode.Mode,
) (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: "write named zettel to index",
		}

		return
	}

	if mode.Contains(objekte_mode.ModeAddToBestandsaufnahme) {
		// TAI must be set before calculating objekte sha
		if mode.Contains(objekte_mode.ModeUpdateTai) {
			t.SetTai(kennung.NowTai())
		}
	}

	if err = t.CalculateObjekteSha(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mode.Contains(objekte_mode.ModeAddToBestandsaufnahme) {
		s.CommitTransacted(t)
	}

	if err = s.GetVerzeichnisse().ExistsOneSha(
		&t.Metadatei.Sha,
	); err == collections.ErrExists {
		return
	}

	err = nil

	g := gattung.Must(t.Kennung.GetGattung())

	switch g {
	case gattung.Konfig:
		if err = s.StoreUtil.GetKonfig().SetTransacted(
			t,
			s.GetAkten().GetKonfigV0(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Kasten, gattung.Typ, gattung.Etikett:
		if err = s.StoreUtil.GetKonfig().AddTransacted(t, s.GetAkten()); err != nil {
			err = errors.Wrap(err)
			return
		}

	case gattung.Zettel:

		if err = s.GetKonfig().ApplyToSku(t); err != nil {
			err = errors.Wrap(err)
			return
		}

		// if err = s.CalculateAndSetShaTransacted(tz); err != nil {
		// 	err = errors.Wrap(err)
		// 	return
		// }

		if err = s.StoreUtil.GetKennungIndex().AddHinweis(&t.Kennung); err != nil {
			if errors.Is(err, hinweisen.ErrDoesNotExist{}) {
				errors.Log().Printf("kennung does not contain value: %s", err)
				err = nil
			} else {
				err = errors.Wrapf(err, "failed to write zettel to index: %s", t)
				return
			}
		}

	default:
		err = gattung.MakeErrUnsupportedGattung(g)
		return
	}

	if err = s.GetVerzeichnisse().Add(
		t,
		t.GetKennung().String(),
		mode,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) handleNew(
	t *sku.Transacted,
	updateType objekte_mode.Mode,
) (err error) {
	if err = s.handleNewOrUpdated(t, updateType); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return s.New(t)
}

func (s *Store) handleUpdated(
	t *sku.Transacted,
	updateType objekte_mode.Mode,
) (err error) {
	if err = s.handleNewOrUpdated(t, updateType); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return s.Updated(t)
}

func (s *Store) handleUnchanged(
	t *sku.Transacted,
) (err error) {
	return s.Unchanged(t)
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

		if sk, err = s.GetVerzeichnisse().ReadOneKennung(h); err != nil {
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
		err = errors.Wrap(objekte_store.ErrNotFound(k1.String()))
		return
	}

	if err = out.SetFromSkuLike(sk); err != nil {
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
		return s.ReadAllSchwanzen(gattungen.MakeSet(gs...), f)
	}
}

func (s *Store) ReadAllSchwanzen(
	gs gattungen.Set,
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
			f,
		),
	)
}

func (s *Store) ReadAll(
	gs gattungen.Set,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	return s.ReadAllGattungenFromVerzeichnisse(gs, f)
}
