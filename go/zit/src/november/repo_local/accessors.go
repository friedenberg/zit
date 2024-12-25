package repo_local

import (
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (u *Repo) GetTime() time.Time {
	return time.Now()
}

func (u *Repo) GetConfig() *config.Compiled {
	return &u.config
}

func (u *Repo) GetDormantIndex() *dormant_index.Index {
	return &u.dormantIndex
}

func (u *Repo) GetDirectoryLayout() dir_layout.DirLayout {
	return u.dirLayout
}

func (u *Repo) GetBlobStore() interfaces.BlobStore {
	return u.GetDirectoryLayout()
}

func (u *Repo) GetStore() *store.Store {
	return &u.store
}

func (u *Repo) GetExternalLikePoolForRepoId(
	repoId ids.RepoId,
) (of sku.ObjectFactory) {
	return
}

func (u *Repo) GetFileEncoder() store_fs.FileEncoder {
	return u.fileEncoder
}
