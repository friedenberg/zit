package store

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/query"
)

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

		sk = s.GetKonfig().GetApproximatedTyp(k).ActualOrNil()

	case gattung.Etikett:
		var e kennung.Etikett

		if err = e.Set(k1.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		ok := false
		sk, ok = s.GetKonfig().GetEtikett(e)

		if !ok {
			sk = nil
		}

	case gattung.Kasten:
		var k kennung.Kasten

		if err = k.Set(k.String()); err != nil {
			err = errors.Wrap(err)
			return
		}

		sk = s.GetKonfig().GetKasten(k)

	case gattung.Konfig:
		sk = &s.GetKonfig().Sku

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

	return
}

func (s *Store) ReadAllGattungFromBestandsaufnahme(
	g gattung.Gattung,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	eachSku := func(_, sk *sku.Transacted) (err error) {
		if sk.GetGattung() != g {
			return
		}

		if err = f(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.GetBestandsaufnahmeStore().ReadAllSkus(eachSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadAllGattungenFromBestandsaufnahme(
	g kennung.Gattung,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if g.IsEmpty() {
		return
	}

	eachSku := func(besty, sk *sku.Transacted) (err error) {
		if !g.ContainsOneOf(sk.GetGattung()) {
			return
		}

		if err = f(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.GetBestandsaufnahmeStore().ReadAllSkus(eachSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadAllGattungFromVerzeichnisse(
	qg *query.Group,
	g gattung.Gattung,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	eachSku := func(sk *sku.Transacted) (err error) {
		if sk.GetGattung() != g {
			return
		}

		if err = f(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.verzeichnisse.ReadQuery(
		qg,
		eachSku,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadAllGattungenFromVerzeichnisse(
	qg *query.Group,
	g kennung.Gattung,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if g.IsEmpty() {
		return
	}

	eachSku := func(sk *sku.Transacted) (err error) {
		if !g.ContainsOneOf(sk.GetGattung()) {
			return
		}

		if err = f(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.verzeichnisse.ReadQuery(
		qg,
		eachSku,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
