package store_browser

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) DeleteExternalLike(el sku.ExternalLike) (err error) {
	e := el.(*sku.Transacted)

	var item Item

	if err = item.ReadFromExternal(e); err != nil {
		err = errors.Wrap(err)
		return
	}

	item.ExternalId = e.GetSku().ObjectId.String()

	s.deleted[item.Url.URL] = append(s.deleted[item.Url.URL], transactedWithItem{
		Transacted: e.CloneTransacted(),
		Item:       item,
	})

	return
}
