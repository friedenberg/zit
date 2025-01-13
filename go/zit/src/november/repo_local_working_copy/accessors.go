package repo_local_working_copy

import (
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/mike/config"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (u *Repo) GetEnv() *env.Env {
	return u.Env
}

func (u *Repo) GetTime() time.Time {
	return time.Now()
}

func (u *Repo) GetConfig() *config.Compiled {
	return &u.config
}

func (u *Repo) GetDormantIndex() *dormant_index.Index {
	return &u.dormantIndex
}

func (u *Repo) GetRepoLayout() repo_layout.Layout {
	return u.layout
}

func (u *Repo) GetBlobStore() interfaces.BlobStore {
	return u.GetRepoLayout()
}

func (u *Repo) GetInventoryListStore() sku.InventoryListStore {
	return u.GetStore().GetInventoryListStore()
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
