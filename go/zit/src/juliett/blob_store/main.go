package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
)

type VersionedStores struct {
	InventoryList InventoryStore
	Repo          RepoStore
	Config        ConfigStore
	Type          TypeStore
	Tag           TagStore
}

func Make(
	dirLayout dir_layout.DirLayout,
	luaVMPoolBuilder *lua.VMPoolBuilder,
	objectFormat object_inventory_format.Format,
	boxFormat *box_format.BoxTransacted,
) *VersionedStores {
	return &VersionedStores{
		InventoryList: MakeInventoryStore(dirLayout, objectFormat, boxFormat),
		Tag:           MakeTagStore(dirLayout, luaVMPoolBuilder),
		Repo:          MakeRepoStore(dirLayout),
		Config:        MakeConfigStore(dirLayout),
		Type:          MakeTypeStore(dirLayout),
	}
}

func (a *VersionedStores) GetTypeV1() Store[type_blobs.TomlV1, *type_blobs.TomlV1] {
	return a.Type.toml_v1
}

func (a *VersionedStores) GetType() TypeStore {
	return a.Type
}

func (a *VersionedStores) GetConfig() ConfigStore {
	return a.Config
}

func (a *VersionedStores) GetTag() TagStore {
	return a.Tag
}

func (a *VersionedStores) GetInventoryList() InventoryStore {
	return a.InventoryList
}
