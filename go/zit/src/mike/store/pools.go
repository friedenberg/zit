package store

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/chrome"
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
