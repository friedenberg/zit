package repo_signing

import (
	"crypto"

	"code.linenisgreat.com/zit/go/zit/src/bravo/bech32"
)

type TomlPublicKeyV0 struct {
	PublicKey bech32.Value `toml:"public-key,omitempty"`
}

func (b TomlPublicKeyV0) GetPublicKey() PublicKey {
	return b.PublicKey.Data
}

func (b *TomlPublicKeyV0) SetPublicKey(key crypto.PublicKey) {
	b.PublicKey.HRP = "zit-repo-public_key-v0"
	b.PublicKey.Data = key.(PublicKey)
}
