package store_util

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type mutators interface {
	CommitTransacted(*sku.Transacted) error
}

func (s *Store) CommitTransacted(t *sku.Transacted) (err error) {
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

func (s *Store) AddTypToIndex(t *kennung.Typ) (err error) {
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
