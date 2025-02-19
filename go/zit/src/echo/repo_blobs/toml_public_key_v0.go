package repo_blobs

import (
	"crypto"
	"crypto/ed25519"

	"code.linenisgreat.com/zit/go/zit/src/bravo/bech32"
)

type tomlPublicKeyV0 struct {
	PublicKey bech32.Value `toml:"public-key,omitempty"`
}

func (b tomlPublicKeyV0) GetPublicKey() ed25519.PublicKey {
	return b.PublicKey.Data
}

func (b *tomlPublicKeyV0) SetPublicKey(key crypto.PublicKey) {
	b.PublicKey.Data = key.(ed25519.PublicKey)
}
