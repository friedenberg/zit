package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type reader interface {
	ReadAllGattungFromBestandsaufnahme(
		g gattung.Gattung,
		f schnittstellen.FuncIter[*sku.Transacted],
	) (err error)

	ReadAllGattungenFromBestandsaufnahme(
		g gattungen.Set,
		f schnittstellen.FuncIter[*sku.Transacted],
	) (err error)

	ReadAllGattungFromVerzeichnisse(
		g gattung.Gattung,
		f schnittstellen.FuncIter[*sku.Transacted],
	) (err error)

	ReadAllGattungenFromVerzeichnisse(
		g gattungen.Set,
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
	g gattungen.Set,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if g.Len() == 0 {
		return
	}

	eachSku := func(besty, sk *sku.Transacted) (err error) {
		if !g.ContainsKey(sk.GetGattung().GetGattungString()) {
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
	g gattungen.Set,
	f schnittstellen.FuncIter[*sku.Transacted],
) (err error) {
	if g.Len() == 0 {
		return
	}

	eachSku := func(sk *sku.Transacted) (err error) {
		if !g.ContainsKey(sk.GetGattung().GetGattungString()) {
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
