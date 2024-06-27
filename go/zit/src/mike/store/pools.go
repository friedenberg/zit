package store

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
)

func (s *Store) PutCheckedOutLike(col sku.CheckedOutLike) {
	switch col.GetKasten().GetKastenString() {
	default:
		cofs := col.(*store_fs.CheckedOut)
		store_fs.GetCheckedOutPool().Put(cofs)
	}
}
