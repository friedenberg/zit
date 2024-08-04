package browser

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) DeleteCheckout(col sku.CheckedOutLike) (err error) {
	coc := col.(*CheckedOut)

	var u *url.URL

	if u, err = coc.External.item.GetUrl(); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.removed[*u] = struct{}{}

	if err = s.itemDeletedStringFormatWriter(coc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
