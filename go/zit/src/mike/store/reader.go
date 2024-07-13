package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO add support for cwd and sigil
// TODO simplify
func (s *Store) ReadOneInto(
	k1 interfaces.StringerGenreGetter,
	out *sku.Transacted,
) (err error) {
	var sk *sku.Transacted

	switch k1.GetGenre() {
	case genres.Zettel:
		var h *ids.ZettelId

		if h, err = s.GetAbbrStore().Hinweis().ExpandString(
			k1.String(),
		); err == nil {
			k1 = h
		} else {
			err = nil
		}

		if sk, err = s.ReadOneKennung(k1); err != nil {
			err = errors.Wrap(err)
			return
		}

	case genres.Type, genres.Tag, genres.Repo:
		if sk, err = s.ReadOneKennung(k1); err != nil {
			err = errors.Wrap(err)
			return
		}

	// case gattung.Typ:
	// 	var k kennung.Typ

	// 	if err = k.Set(k1.String()); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	sk = s.GetKonfig().GetApproximatedTyp(k).ActualOrNil()

	// case gattung.Etikett:
	// 	var e kennung.Etikett

	// 	if err = e.Set(k1.String()); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	ok := false
	// 	sk, ok = s.GetKonfig().GetEtikett(e)

	// 	if !ok {
	// 		sk = nil
	// 	}

	// case gattung.Kasten:
	// 	var k kennung.Kasten

	// 	if err = k.Set(k.String()); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}

	// 	sk = s.GetKonfig().GetKasten(k)

	case genres.Config:
		sk = &s.GetKonfig().Sku

		if sk.GetTai().IsEmpty() {
			sk = nil
		}

	default:
		err = genres.MakeErrUnsupportedGattung(k1)
		return
	}

	if sk == nil {
		err = collections.MakeErrNotFound(k1)
		return
	}

	if err = out.SetFromSkuLike(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
