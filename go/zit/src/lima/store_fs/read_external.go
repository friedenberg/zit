package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (store *Store) ReadExternalLikeFromObjectIdLike(
	commitOptions sku.CommitOptions,
	objectIdMaybeExternal interfaces.Stringer,
	internal *sku.Transacted,
) (external sku.ExternalLike, err error) {
	var items []*sku.FSItem

	oidString := objectIdMaybeExternal.String()
	_, isExternal := objectIdMaybeExternal.(ids.ExternalObjectIdLike)

	if !isExternal {
		oidString = store.keyForObjectIdString(oidString)
	}

	if items, err = store.GetFSItemsForString(
		store.envRepo.GetCwd(),
		oidString,
		true,
	); err != nil {
		err = errors.Wrapf(err, "ObjectIdString: %q", oidString)
		return
	}

	switch len(items) {
	case 0:
		if !isExternal {
			external = sku.GetTransactedPool().Get()

			var objectId ids.ObjectId

			if err = objectId.Set(oidString); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = store.externalStoreSupplies.ReadOneInto(
				&objectId,
				external.GetSku(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return

	case 1:
		break

	default:
		err = errors.Errorf(
			"more than one FSItem (%q) matches object id (%q).",
			items,
			objectIdMaybeExternal,
		)

		return
	}

	item := items[0]

	if external, err = store.ReadExternalFromItem(
		commitOptions,
		item,
		internal,
	); err != nil {
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
