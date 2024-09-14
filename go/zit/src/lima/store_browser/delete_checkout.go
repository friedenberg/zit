package store_browser

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) DeleteExternalLike(el sku.ExternalLike) (err error) {
	e := el.(*External)

	bi := e.Item
	bi.ExternalId = e.GetSku().ObjectId.String()
	s.deleted[e.Item.Url.URL] = append(s.deleted[e.Item.Url.URL], bi)

	return
}
