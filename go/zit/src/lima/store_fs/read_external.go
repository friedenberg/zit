package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadExternalLikeFromObjectId(
	o sku.CommitOptions,
	oid interfaces.ObjectId,
	internal *sku.Transacted,
) (external sku.ExternalLike, err error) {
	item, ok := s.Get(oid)

	if !ok {
		return
	}

	if external, err = s.ReadExternalFromItem(o, item, internal); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// Given a sku and an FSItem, return the overlayed external variant. Internal
// can be nil and then only the external data is used.
func (s *Store) ReadExternalFromItem(
	o sku.CommitOptions,
	item *sku.FSItem,
	internal *sku.Transacted,
) (external *sku.Transacted, err error) {
	external = GetExternalPool().Get()

	if err = s.HydrateExternalFromItem(o, item, internal, external); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
