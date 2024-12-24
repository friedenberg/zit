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

	item.ExternalId = external.GetSkuExternal().GetExternalObjectId().String()

	s.deleted[item.Url.Url()] = append(s.deleted[item.Url.Url()], checkedOutWithItem{
		CheckedOut: co.Clone(),
		Item:       item,
	})

	return
}
