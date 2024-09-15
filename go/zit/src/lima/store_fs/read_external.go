package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadExternalLikeFromObjectId(
	o sku.CommitOptions,
	k1 interfaces.ObjectId,
	t *sku.Transacted,
) (e sku.ExternalLike, err error) {
	k, ok := s.Get(k1)

	if !ok {
		return
	}

	if e, err = s.ReadExternalFromObjectIdFDPair(o, k, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadExternalFromObjectIdFDPair(
	o sku.CommitOptions,
	em *Item,
	t *sku.Transacted,
) (e *External, err error) {
	e = GetExternalPool().Get()

	if err = s.ReadIntoExternalFromItem(o, em, t, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadIntoExternalFromItem(
	o sku.CommitOptions,
	i *Item,
	t *sku.Transacted,
	e *External,
) (err error) {
	o.Del(objekte_mode.ModeApplyProto)

	if err = s.ReadOneExternalInto(&o, i, t, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.externalStoreSupplies.FuncCommit(
		&e.Transacted,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
