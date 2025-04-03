package local_working_copy

import (
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/dormant_index"
	"code.linenisgreat.com/zit/go/zit/src/kilo/env_workspace"
	"code.linenisgreat.com/zit/go/zit/src/lima/env_lua"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/mike/store_config"
)

func (u *Repo) GetEnv() env_ui.Env {
	return u
}

func (u *Repo) GetImmutableConfigPublic() config_immutable_io.ConfigLoadedPublic {
	return u.GetEnvRepo().GetConfigPublic()
}

func (u *Repo) GetImmutableConfigPrivate() config_immutable_io.ConfigLoadedPrivate {
	return u.GetEnvRepo().GetConfigPrivate()
}

func (repo *Repo) GetEnvLocal() env_local.Env {
	return repo
}

func (repo *Repo) GetEnvWorkspace() env_workspace.Env {
	return repo.envWorkspace
}

func (u *Repo) GetEnvLua() env_lua.Env {
	return u.envLua
}

func (u *Repo) GetTime() time.Time {
	return time.Now()
}

func (u *Repo) GetConfig() store_config.Store {
	return u.config
}

func (u *Repo) GetDormantIndex() *dormant_index.Index {
	return &u.dormantIndex
}

func (u *Repo) GetEnvRepo() env_repo.Env {
	return u.envRepo
}

func (u *Repo) GetTypedInventoryListBlobStore() typed_blob_store.InventoryList {
	return u.typedBlobStore.InventoryList
}

func (u *Repo) GetBlobStore() interfaces.BlobStore {
	return u.GetEnvRepo()
}

func (repo *Repo) GetObjectStore() sku.ObjectStore {
	return &repo.store
}

func (u *Repo) GetInventoryListStore() sku.InventoryListStore {
	return u.GetStore().GetInventoryListStore()
}

func (u *Repo) GetStore() *store.Store {
	return &u.store
}

func (u *Repo) GetFileEncoder() store_fs.FileEncoder {
	return u.fileEncoder
}
