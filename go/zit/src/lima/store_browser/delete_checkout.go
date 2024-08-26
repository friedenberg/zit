package store_browser

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) DeleteCheckout(col sku.CheckedOutLike) (err error) {
	coc := col.(*CheckedOut)

	var u *url.URL

	if u, err = coc.External.browserItem.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

  bi := coc.External.browserItem
  bi.ExternalId = coc.GetSku().ObjectId.String()
	s.removed[*u] = append(s.removed[*u], bi)

	return
}
