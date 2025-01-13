package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
)

type VersionedStores struct {
	InventoryList InventoryList
	Repo          RepoStore
	Config        Config
	Type          Type
	Tag           Tag
}

func Make(
	dirLayout repo_layout.Layout,
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

func (a *VersionedStores) GetType() Type {
	return a.Type
}

func (a *VersionedStores) GetConfig() Config {
	return a.Config
}

func (a *VersionedStores) GetTag() Tag {
	return a.Tag
}

func (a *VersionedStores) GetInventoryList() InventoryList {
	return a.InventoryList
}
