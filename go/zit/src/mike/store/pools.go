package store

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/browser"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
)

func (s *Store) PutCheckedOutLike(col sku.CheckedOutLike) {
	switch col.GetRepoId().GetRepoIdString() {
	// TODO make generic?
	case "browser":
		browser.GetCheckedOutPool().Put(col.(*browser.CheckedOut))

	default:
		cofs := col.(*store_fs.CheckedOut)
		store_fs.GetCheckedOutPool().Put(cofs)
	}
}
