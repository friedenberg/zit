package store_util

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type reader interface {
	ReadAllGattungFromBestandsaufnahme(
		g gattung.Gattung,
		f schnittstellen.FuncIter[*sku.Transacted],
	) (err error)

	ReadAllGattungenFromBestandsaufnahme(
		g kennung.Gattung,
		f schnittstellen.FuncIter[*sku.Transacted],
	) (err error)

	ReadAllGattungFromVerzeichnisse(
		g gattung.Gattung,
		f schnittstellen.FuncIter[*sku.Transacted],
	) (err error)

	ReadAllGattungenFromVerzeichnisse(
		g kennung.Gattung,
		f schnittstellen.FuncIter[*sku.Transacted],
	) (err error)
}

func (s *common) ReadAllGattungFromBestandsaufnahme(
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

func (s *common) ReadAllGattungenFromBestandsaufnahme(
	g kennung.Gattung,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if g.IsEmpty() {
		return
	}

	eachSku := func(besty, sk *sku.Transacted) (err error) {
		if !g.Contains(sk.GetGattung()) {
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

func (s *common) ReadAllGattungFromVerzeichnisse(
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

	if err = s.verzeichnisse.ReadAll(eachSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *common) ReadAllGattungenFromVerzeichnisse(
	g kennung.Gattung,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if g.IsEmpty() {
		return
	}

	eachSku := func(sk *sku.Transacted) (err error) {
		if !g.Contains(sk.GetGattung()) {
			return
		}

		if err = f(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.verzeichnisse.ReadAll(eachSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
