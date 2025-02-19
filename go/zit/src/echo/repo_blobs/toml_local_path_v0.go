package repo_blobs

import (
	"crypto/ed25519"

	"code.linenisgreat.com/zit/go/zit/src/bravo/bech32"
)

type TomlLocalPathV0 struct {
	PublicKey bech32.Value `toml:"public-key,omitempty"`
	Path      string       `toml:"path"`
}

func (b TomlLocalPathV0) GetRepoBlob() Blob {
	return b
}

func (b TomlLocalPathV0) GetPublicKey() ed25519.PublicKey {
	return b.PublicKey.Data
}

func (a *TomlLocalPathV0) Reset() {
	a.Path = ""
}

func (a *TomlLocalPathV0) ResetWith(b TomlLocalPathV0) {
	a.Path = b.Path
}

func (a TomlLocalPathV0) Equals(b TomlLocalPathV0) bool {
	if a.Path != b.Path {
		return false
	}

	return true
}
