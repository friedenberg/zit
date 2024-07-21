package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

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
