package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadOneObjectId(
	k interfaces.ObjectId,
) (sk *sku.Transacted, err error) {
	sk = sku.GetTransactedPool().Get()

	if err = s.GetStreamIndex().ReadOneObjectId(k.String(), sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReaderFor(sh *sha.Sha) (rc sha.ReadCloser, err error) {
	if rc, err = s.fs_home.BlobReaderFrom(
		sh,
		s.fs_home.DirVerzeichnisseMetadataObjectIdParent(),
	); err != nil {
		if errors.IsNotExist(err) {
			err = collections.MakeErrNotFound(sh)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

// TODO add support for cwd and sigil
// TODO simplify
func (s *Store) ReadOneInto(
	k1 interfaces.ObjectId,
	out *sku.Transacted,
) (err error) {
	var sk *sku.Transacted

	switch k1.GetGenre() {
	case genres.Zettel:
		var h *ids.ZettelId

		if h, err = s.GetAbbrStore().ZettelId().ExpandString(
			k1.String(),
		); err == nil {
			k1 = h
		} else {
			err = nil
		}

		if sk, err = s.ReadOneObjectId(k1); err != nil {
			err = errors.Wrap(err)
			return
		}

	case genres.Type, genres.Tag, genres.Repo:
		if sk, err = s.ReadOneObjectId(k1); err != nil {
			err = errors.Wrap(err)
			return
		}

	case genres.Config:
		sk = &s.GetKonfig().Sku

		if sk.GetTai().IsEmpty() {
			ui.Err().Print("config tai is empty")
		}

	case genres.Blob:
		var oid ids.ObjectId

		if err = oid.SetWithGenre(k1.String(), genres.Zettel); err != nil {
			err = collections.MakeErrNotFound(k1)
			return
		}

		if sk, err = s.ReadOneObjectId(k1); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = genres.MakeErrUnsupportedGenre(k1)
		return
	}

	if sk == nil {
		err = collections.MakeErrNotFound(k1)
		return
	}

	sku.TransactedResetter.ResetWith(out, sk)

	return
}
