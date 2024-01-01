package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type mutators interface {
	AddVerzeichnisse(*sku.Transacted) error
	CommitTransacted(*sku.Transacted) error
	CommitUpdatedTransacted(*sku.Transacted) error
}

func (s *common) AddVerzeichnisse(t *sku.Transacted) (err error) {
	if err = s.verzeichnisse.Add(
		t,
		t.GetKennung().String(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *common) CommitUpdatedTransacted(
	t *sku.Transacted,
) (err error) {
  // log.Debug().Printf("%s", t.StringKennungTai())
	ta := kennung.NowTai()

	if t.GetTai().IsZero() {
		t.SetTai(ta)
	}

	return s.CommitTransacted(t)
}

func (s *common) CommitTransacted(t *sku.Transacted) (err error) {
	sk := sku.GetTransactedPool().Get()

	if err = s.konfig.AddTransacted(
		t,
		s.GetAkten(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = sk.SetFromSkuLike(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.bestandsaufnahmeAkte.Skus.Add(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *common) AddTypToIndex(t *kennung.Typ) (err error) {
	if t == nil {
		return
	}

	if t.IsEmpty() {
		return
	}

	if err = s.typenIndex.StoreOne(*t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
