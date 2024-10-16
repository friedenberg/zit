package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
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

	if e, err = s.ReadExternalFromItem(o, k, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadExternalFromItem(
	o sku.CommitOptions,
	i *Item,
	t *sku.Transacted,
) (e *sku.Transacted, err error) {
	e = GetExternalPool().Get()

	if err = s.readIntoExternalFromItem(o, i, t, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) readIntoExternalFromItem(
	o sku.CommitOptions,
	i *Item,
	t *sku.Transacted,
	e *sku.Transacted,
) (err error) {
	o.Del(object_mode.ModeApplyProto)

	if err = s.readOneExternalInto(&o, i, t, e); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.externalStoreSupplies.FuncCommit(
		e,
		o,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
