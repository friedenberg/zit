package typed_blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
)

type Store struct {
	InventoryList InventoryList
	Repo          RepoStore
	Config        Config
	Type          Type
	Tag           Tag
}

func Make(
	envRepo env_repo.Env,
	luaVMPoolBuilder *lua.VMPoolBuilder,
	objectFormat object_inventory_format.Format,
	boxFormat *box_format.BoxTransacted,
) *Store {
	return &Store{
		InventoryList: MakeInventoryStore(envRepo, objectFormat, boxFormat),
		Tag:           MakeTagStore(envRepo, luaVMPoolBuilder),
		Repo:          MakeRepoStore(envRepo),
		Config:        MakeConfigStore(envRepo),
		Type:          MakeTypeStore(envRepo),
	}
}

func (a *Store) GetTypeV1() TypedStore[type_blobs.TomlV1, *type_blobs.TomlV1] {
	return a.Type.toml_v1
}

func (a *Store) GetType() Type {
	return a.Type
}

func (a *Store) GetConfig() Config {
	return a.Config
}

func (a *Store) GetTag() Tag {
	return a.Tag
}

func (a *Store) GetInventoryList() InventoryList {
	return a.InventoryList
}
