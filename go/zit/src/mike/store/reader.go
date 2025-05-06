package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (s *Store) ReadTransactedFromObjectId(
	k1 interfaces.ObjectId,
) (sk1 *sku.Transacted, err error) {
	sk1 = sku.GetTransactedPool().Get()

	if err = s.ReadOneInto(k1, sk1); err != nil {
		if collections.IsErrNotFound(err) {
			sku.GetTransactedPool().Put(sk1)
			sk1 = nil
		}

		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadOneObjectId(
	k interfaces.ObjectId,
) (sk *sku.Transacted, err error) {
	sk = sku.GetTransactedPool().Get()

	if err = s.GetStreamIndex().ReadOneObjectId(k, sk); err != nil {
		err = errors.WrapExcept(err, collections.ErrNotFound)
		return
	}

	return
}

// TODO add support for cwd and sigil
// TODO simplify
func (s *Store) ReadOneInto(
	objectId interfaces.ObjectId,
	out *sku.Transacted,
) (err error) {
	var sk *sku.Transacted

	switch objectId.GetGenre() {
	case genres.Zettel:
		var zettelId *ids.ZettelId

		if zettelId, err = s.GetAbbrStore().ZettelId().ExpandString(
			objectId.String(),
		); err == nil {
			objectId = zettelId
		} else {
			err = nil
		}

		if sk, err = s.ReadOneObjectId(objectId); err != nil {
			err = errors.Wrap(err)
			return
		}

	case genres.Type, genres.Tag, genres.Repo, genres.InventoryList:
		if sk, err = s.ReadOneObjectId(objectId); err != nil {
			err = errors.Wrap(err)
			return
		}

	case genres.Config:
		sk = s.GetConfig().GetSku()

		if sk.GetTai().IsEmpty() {
			ui.Err().Print("config tai is empty")
		}

	case genres.Blob:
		var oid ids.ObjectId

		if err = oid.SetWithIdLike(objectId); err != nil {
			err = collections.MakeErrNotFound(objectId)
			return
		}

		if sk, err = s.ReadOneObjectId(objectId); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = genres.MakeErrUnsupportedGenre(objectId)
		return
	}

	if sk == nil {
		err = collections.MakeErrNotFound(objectId)
		return
	}

	sku.TransactedResetter.ResetWith(out, sk)

	return
}
