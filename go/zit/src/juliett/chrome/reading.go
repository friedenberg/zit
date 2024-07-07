package chrome

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO decide how this should behave
func (s *Store) UpdateTransacted(sk *sku.Transacted) (err error) {
	if !sk.GetTyp().Equals(s.typ) {
		return
	}

	var uSku *url.URL

	if uSku, err = s.getUrl(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	_, ok := s.urls[*uSku]

	if !ok {
		return
	}

	return
}
