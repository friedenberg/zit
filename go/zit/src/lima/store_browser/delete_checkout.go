package store_browser

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) DeleteExternalLike(el sku.ExternalLike) (err error) {
	ui.Debug().Print(el)
	e := el.(*External)

	var u *url.URL

	if u, err = e.Item.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

	bi := e.Item
	bi.ExternalId = e.GetSku().ObjectId.String()
	s.deleted[*u] = append(s.deleted[*u], bi)

	return
}
