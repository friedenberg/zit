package store

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/chrome"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
)

func (s *Store) PutCheckedOutLike(col sku.CheckedOutLike) {
	switch col.GetRepoId().GetRepoIdString() {
	case "chrome":
		chrome.GetCheckedOutPool().Put(col.(*chrome.CheckedOut))

	default:
		cofs := col.(*store_fs.CheckedOut)
		store_fs.GetCheckedOutPool().Put(cofs)
	}
}
