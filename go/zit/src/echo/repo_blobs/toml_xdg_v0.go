package repo_blobs

import (
	"crypto/ed25519"

	"code.linenisgreat.com/zit/go/zit/src/bravo/bech32"
	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
)

type TomlXDGV0 struct {
	PublicKey bech32.Value `toml:"public-key,omitempty"`
	Data      string       `toml:"data"`
	Config    string       `toml:"config"`
	State     string       `toml:"state"`
	Cache     string       `toml:"cache"`
	Runtime   string       `toml:"runtime"`
}

func TomlXDGV0FromXDG(xdg xdg.XDG) TomlXDGV0 {
	return TomlXDGV0{
		Data:    xdg.Data,
		Config:  xdg.Config,
		State:   xdg.State,
		Cache:   xdg.Cache,
		Runtime: xdg.Runtime,
	}
}

func (b TomlXDGV0) GetRepoBlob() Blob {
	return b
}

func (b TomlXDGV0) GetPublicKey() ed25519.PublicKey {
	return b.PublicKey.Data
}
