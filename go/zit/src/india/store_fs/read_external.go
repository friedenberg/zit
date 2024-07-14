package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadTransactedFromKennung(
	o sku.ObjectOptions,
	k1 interfaces.ObjectId,
	t *sku.Transacted,
) (e *External, err error) {
	k, ok := s.Get(k1)

	if !ok {
		return
	}

	if e, err = s.ReadExternalFromKennungFDPair(o, k, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadExternalFromKennungFDPair(
	o sku.ObjectOptions,
	em *KennungFDPair,
	t *sku.Transacted,
) (e *External, err error) {
	e = GetExternalPool().Get()

	if err = s.ReadIntoExternalFromKennungFDPair(o, em, t, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadIntoExternalFromKennungFDPair(
	o sku.ObjectOptions,
	em *KennungFDPair,
	t *sku.Transacted,
	e *External,
) (err error) {
	o.Del(objekte_mode.ModeApplyProto)

	if err = s.ReadOneExternalInto(
		&o,
		em,
		t,
		e,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.externalStoreInfo.FuncCommit(
		&e.Transacted,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
