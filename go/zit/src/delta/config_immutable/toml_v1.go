package config_immutable

import (
	"crypto/ed25519"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/bravo/bech32"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type TomlV1 struct {
	PrivateKey   bech32.Value    `toml:"private-key,omitempty"`
	StoreVersion StoreVersion    `toml:"store-version"`
	RepoType     repo_type.Type  `toml:"repo-type"`
	RepoId       ids.RepoId      `toml:"id"`
	BlobStore    BlobStoreTomlV1 `toml:"blob-store"`
}

func (k *TomlV1) SetFlagSet(f *flag.FlagSet) {
	k.BlobStore.SetFlagSet(f)
	k.RepoType = repo_type.TypeWorkingCopy
	f.Var(&k.RepoType, "repo-type", "")
}

func (k *TomlV1) GetImmutableConfig() Config {
	return k
}

func (k *TomlV1) GetBlobStoreConfigImmutable() interfaces.BlobStoreConfigImmutable {
	return &k.BlobStore
}

func (k *TomlV1) GetPrivateKey() ed25519.PrivateKey {
	return ed25519.NewKeyFromSeed(k.PrivateKey.Data)
}

func (k *TomlV1) GetStoreVersion() interfaces.StoreVersion {
	return k.StoreVersion
}

func (k TomlV1) GetRepoType() repo_type.Type {
	return k.RepoType
}

func (k TomlV1) GetRepoId() ids.RepoId {
	return k.RepoId
}
