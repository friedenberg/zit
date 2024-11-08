package store_fs

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) UpdateTransacted(internal *sku.Transacted) (err error) {
	item, ok := s.Get(&internal.ObjectId)

	if !ok {
		return
	}

	var external *sku.Transacted

	if external, err = s.ReadExternalFromItem(
		sku.CommitOptions{
			Mode: object_mode.ModeUpdateTai,
		},
		item,
		internal,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sku.Resetter.ResetWith(internal, external)
	sku.GetTransactedPool().Put(external)

	return
}

func (s *Store) ReadOneExternalObjectReader(
	r io.Reader,
	external *sku.Transacted,
) (err error) {
	if _, err = s.metadataTextParser.ParseMetadata(r, external); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
