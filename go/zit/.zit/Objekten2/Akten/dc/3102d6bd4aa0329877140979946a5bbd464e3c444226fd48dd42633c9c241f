package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
)

func (s *Store) validate(
	el sku.ExternalLike, mutter *sku.Transacted,
	o sku.CommitOptions,
) (err error) {
	if o.DontValidate {
		return
	}

	switch el.GetSku().GetGenre() {
	case genres.Type:
		tipe := el.GetSku().GetType()

		var commonBlob type_blobs.Blob

		if commonBlob, _, err = s.GetBlobStore().GetType().ParseTypedBlob(
			tipe,
			el.GetSku().GetBlobSha(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		defer s.GetBlobStore().GetType().PutTypedBlob(tipe, commonBlob)
	}

	return
}
