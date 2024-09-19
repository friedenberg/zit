package store

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_browser"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
)

func (s *Store) PutCheckedOutLike(col sku.CheckedOutLike) {
	switch col.GetSkuExternalLike().GetRepoId().GetRepoIdString() {
	// TODO make generic?
	case "browser":
		store_browser.GetCheckedOutPool().Put(col.(*store_browser.CheckedOut))

	default:
		cofs := col.(*sku.CheckedOut)
		store_fs.GetCheckedOutPool().Put(cofs)
	}
}
