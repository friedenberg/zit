package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/lua"
	"code.linenisgreat.com/zit/go/zit/src/delta/tag_blobs"
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_blobs"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/mutable_config_blobs"
)

// TODO switch to interfaces instead of structs
type VersionedStores struct {
	repo_v0      Store[repo_blobs.V0, *repo_blobs.V0]
	config_store ConfigStore
	type_store   TypeStore
	tag_store    TagStore
}

func Make(
	dirLayout dir_layout.DirLayout,
	luaVMPoolBuilder *lua.VMPoolBuilder,
) *VersionedStores {
	return &VersionedStores{
		tag_store: MakeTagStore(dirLayout, luaVMPoolBuilder),
		repo_v0: MakeBlobStore(
			dirLayout,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[repo_blobs.V0](
					dirLayout,
				),
				ParsedBlobTomlFormatter[repo_blobs.V0, *repo_blobs.V0]{},
				dirLayout,
			),
			func(a *repo_blobs.V0) {
				a.Reset()
			},
		),
		config_store: MakeConfigStore(dirLayout),
		type_store:   MakeTypeStore(dirLayout),
	}
}

func (a *VersionedStores) GetTagTomlV0() Store[tag_blobs.V0, *tag_blobs.V0] {
	return a.tag_store.toml_v0
}

func (a *VersionedStores) GetTagTomlV1() Store[tag_blobs.TomlV1, *tag_blobs.TomlV1] {
	return a.tag_store.toml_v1
}

func (a *VersionedStores) GetRepoV0() Store[repo_blobs.V0, *repo_blobs.V0] {
	return a.repo_v0
}

func (a *VersionedStores) GetConfigV0() Store[mutable_config_blobs.V0, *mutable_config_blobs.V0] {
	return a.config_store.config_toml_v0
}

func (a *VersionedStores) GetTypeV0() Store[type_blobs.V0, *type_blobs.V0] {
	return a.type_store.type_toml_v0
}

func (a *VersionedStores) GetTypeV1() Store[type_blobs.TomlV1, *type_blobs.TomlV1] {
	return a.type_store.type_toml_v1
}

func (a *VersionedStores) GetConfig() ConfigStore {
	return a.config_store
}

func (a *VersionedStores) GetType() TypeStore {
	return a.type_store
}

func (a *VersionedStores) GetTag() TagStore {
	return a.tag_store
}
