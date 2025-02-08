package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (s *Store) ReadExternalLikeFromObjectIdLike(
	commitOptions sku.CommitOptions,
	objectId interfaces.Stringer,
	internal *sku.Transacted,
) (external sku.ExternalLike, err error) {
	var items []*sku.FSItem

	oidString := objectId.String()

	if _, ok := objectId.(ids.ExternalObjectIdLike); !ok {
		oidString = s.keyForObjectIdString(oidString)
	}

	if items, err = s.GetFSItemsForString(
		oidString,
		true,
	); err != nil {
		err = errors.Wrapf(err, "ObjectIdString: %q", oidString)
		return
	}

	switch len(items) {
	case 0:
		return

	case 1:
		break

	default:
		err = errors.Errorf(
			"more than one FSItem (%q) matches object id (%q).",
			items,
			objectId,
		)

		return
	}

	item := items[0]

	if external, err = s.ReadExternalFromItem(
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
