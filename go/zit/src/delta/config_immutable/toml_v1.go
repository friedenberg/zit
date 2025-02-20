package config_immutable

import (
	"crypto/ed25519"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/bravo/bech32"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type TomlV1Common struct {
	StoreVersion StoreVersion    `toml:"store-version"`
	RepoType     repo_type.Type  `toml:"repo-type"`
	RepoId       ids.RepoId      `toml:"id"`
	BlobStore    BlobStoreTomlV1 `toml:"blob-store"`
}

type TomlV1Private struct {
	PrivateKey bech32.Value `toml:"private-key,omitempty"`
	TomlV1Common
}

type TomlV1Public struct {
	PublicKey bech32.Value `toml:"public-key,omitempty"`
	TomlV1Common
}

func (k *TomlV1Common) SetFlagSet(f *flag.FlagSet) {
	k.BlobStore.SetFlagSet(f)
	k.RepoType = repo_type.TypeWorkingCopy
	f.Var(&k.RepoType, "repo-type", "")
}

func (k *TomlV1Private) GetImmutableConfig() ConfigPrivate {
	return k
}

func (k *TomlV1Private) GetImmutableConfigPublic() ConfigPublic {
	return &TomlV1Public{
		TomlV1Common: k.TomlV1Common,
		PublicKey: bech32.Value{
			HRP:  "zit-repo-public_key-v1",
			Data: k.GetPublicKey(),
		},
	}
}

func (k *TomlV1Private) GetPrivateKey() ed25519.PrivateKey {
	return ed25519.NewKeyFromSeed(k.PrivateKey.Data)
}

func (k *TomlV1Private) GetPublicKey() ed25519.PublicKey {
	return k.GetPrivateKey().Public().(ed25519.PublicKey)
}

func (k *TomlV1Public) GetImmutableConfigPublic() ConfigPublic {
	return k
}

func (k TomlV1Public) GetPublicKey() ed25519.PublicKey {
	return k.PublicKey.Data
}

func (k *TomlV1Common) GetBlobStoreConfigImmutable() interfaces.BlobStoreConfigImmutable {
	return &k.BlobStore
}

func (k *TomlV1Common) GetStoreVersion() interfaces.StoreVersion {
	return k.StoreVersion
}

func (k TomlV1Common) GetRepoType() repo_type.Type {
	return k.RepoType
}

func (k TomlV1Common) GetRepoId() ids.RepoId {
	return k.RepoId
}
