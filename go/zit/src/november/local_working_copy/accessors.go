package local_working_copy

import (
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/mike/env_config"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (u *Repo) GetEnv() env_ui.Env {
	return u
}

func (u *Repo) GetTime() time.Time {
	return time.Now()
}

func (u *Repo) GetConfig() env_config.Env {
	return u.config
}

func (u *Repo) GetDormantIndex() *dormant_index.Index {
	return &u.dormantIndex
}

func (u *Repo) GetRepoLayout() env_repo.Env {
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
