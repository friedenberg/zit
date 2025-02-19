package repo_blobs

import (
	"crypto/ed25519"

	"code.linenisgreat.com/zit/go/zit/src/bravo/bech32"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
)

type TomlUriV0 struct {
	PublicKey bech32.Value `toml:"public-key,omitempty"`
	Uri       values.Uri   `toml:"uri"`
}

func (b TomlUriV0) GetRepoBlob() Blob {
	return b
}

func (b TomlUriV0) GetPublicKey() ed25519.PublicKey {
	return b.PublicKey.Data
}

func (b TomlUriV0) GetRepoType() {
}

func (a *TomlUriV0) Reset() {
	a.Uri = values.Uri{}
}

func (a *TomlUriV0) ResetWith(b TomlUriV0) {
	a.Uri = b.Uri
}

func (a TomlUriV0) Equals(b TomlUriV0) bool {
	if a.Uri != b.Uri {
		return false
	}

	return true
}
