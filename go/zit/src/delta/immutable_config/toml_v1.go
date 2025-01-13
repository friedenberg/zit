package immutable_config

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
)

// TODO Split into repo and blob store configs
type TomlV1 struct {
	StoreVersion StoreVersion    `toml:"store-version"`
	RepoType     repo_type.Type  `toml:"repo-type"`
	BlobStore    BlobStoreTomlV1 `toml:"blob-store"`
}

func (k *TomlV1) SetFlagSet(f *flag.FlagSet) {
	k.BlobStore.SetFlagSet(f)
	k.RepoType = repo_type.TypeReadWrite
	f.Var(&k.RepoType, "repo-type", "")
}

func (k *TomlV1) GetImmutableConfig() Config {
	return k
}

func (k *TomlV1) GetBlobStoreImmutableConfig() BlobStoreConfig {
	return &k.BlobStore
}

func (k *TomlV1) GetStoreVersion() interfaces.StoreVersion {
	return k.StoreVersion
}

func (k TomlV1) GetRepoType() repo_type.Type {
	return k.RepoType
}
