package repo_signing

import (
	"crypto"
	"crypto/ed25519"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/bech32"
)

type TomlPrivateKeyV0 struct {
	PrivateKey bech32.Value `toml:"private-key,omitempty"`
}

func (b *TomlPrivateKeyV0) GeneratePrivateKey() (err error) {
	if len(b.PrivateKey.Data) > 0 {
		err = errors.Errorf("private key data already exists, refusing to generate.")
		return
	}

	var privateKey ed25519.PrivateKey

	if _, privateKey, err = ed25519.GenerateKey(nil); err != nil {
		err = errors.Wrap(err)
		return
	}

	b.PrivateKey.Data = privateKey.Seed()
	b.PrivateKey.HRP = "zit-repo-private_key-v1"

	return
}

func (b TomlPrivateKeyV0) GetPrivateKey() ed25519.PrivateKey {
	return ed25519.NewKeyFromSeed(b.PrivateKey.Data)
}

func (b *TomlPrivateKeyV0) SetPrivateKey(key crypto.PrivateKey) {
	b.PrivateKey.HRP = "zit-repo-private_key-v0"
	b.PrivateKey.Data = key.(ed25519.PrivateKey)
}

func (b *TomlPrivateKeyV0) GetPublicKey() TomlPublicKeyV0 {
	pub := bech32.Value{
		HRP:  "zit-repo-public_key-v0",
		Data: b.GetPrivateKey().Public().(ed25519.PublicKey),
	}

	return TomlPublicKeyV0{PublicKey: pub}
}
