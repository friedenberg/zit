package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/india/tag_blobs"
)

// TODO switch to interfaces instead of structs
type VersionedStores struct {
	inventory_list InventoryStore
	repo           RepoStore
	config         ConfigStore
	tipe           TypeStore
	tag            TagStore
}

func Make(
	dirLayout dir_layout.DirLayout,
	luaVMPoolBuilder *lua.VMPoolBuilder,
	objectFormat object_inventory_format.Format,
	boxFormat *box_format.Box,
) *VersionedStores {
	return &VersionedStores{
		inventory_list: MakeInventoryStore(dirLayout, objectFormat, boxFormat),
		tag:            MakeTagStore(dirLayout, luaVMPoolBuilder),
		repo:           MakeRepoStore(dirLayout),
		config:         MakeConfigStore(dirLayout),
		tipe:           MakeTypeStore(dirLayout),
	}
}

func (a *VersionedStores) GetTagTomlV0() Store[tag_blobs.V0, *tag_blobs.V0] {
	return a.tag.toml_v0
}

func (a *VersionedStores) GetTagTomlV1() Store[tag_blobs.TomlV1, *tag_blobs.TomlV1] {
	return a.tag.toml_v1
}

func (a *VersionedStores) GetConfigV0() Store[mutable_config_blobs.V0, *mutable_config_blobs.V0] {
	return a.config.toml_v0
}

func (a *VersionedStores) GetTypeV0() Store[type_blobs.V0, *type_blobs.V0] {
	return a.tipe.toml_v0
}

func (a *VersionedStores) GetTypeV1() Store[type_blobs.TomlV1, *type_blobs.TomlV1] {
	return a.tipe.toml_v1
}

func (a *VersionedStores) GetConfig() ConfigStore {
	return a.config
}

func (a *VersionedStores) GetType() TypeStore {
	return a.tipe
}

func (a *VersionedStores) GetTag() TagStore {
	return a.tag
}

func (a *VersionedStores) GetInventoryList() InventoryStore {
	return a.inventory_list
}
