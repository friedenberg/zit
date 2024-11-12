package store_browser

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) DeleteCheckedOut(co *sku.CheckedOut) (err error) {
	external := co.GetSkuExternal()

	var item Item

	if err = item.ReadFromExternal(external); err != nil {
		err = errors.Wrap(err)
		return
	}

	item.ExternalId = external.GetSku().ObjectId.String()

	s.deleted[item.Url.URL] = append(s.deleted[item.Url.URL], checkedOutWithItem{
		CheckedOut: co.Clone(),
		Item:       item,
	})

	return
}
