package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (s *Store) ReadExternalLikeFromObjectId(
	commitOptions sku.CommitOptions,
	objectId interfaces.Stringer,
	internal *sku.Transacted,
) (external sku.ExternalLike, err error) {
	var results []*sku.FSItem

	oidString := s.keyForObjectIdString(objectId.String())

	if results, err = s.dirItems.getFDsForObjectIdString(
		oidString,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch len(results) {
	case 0:
		return

	case 1:
		break

	default:
		err = errors.Errorf(
			"more than one FSItem (%q) matches object id (%q).",
			results,
			objectId,
		)

		return
	}

	item := results[0]

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
