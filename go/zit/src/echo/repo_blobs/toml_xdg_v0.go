package repo_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
)

type TomlXDGV0 struct {
	tomlPublicKeyV0
	Data    string `toml:"data"`
	Config  string `toml:"config"`
	State   string `toml:"state"`
	Cache   string `toml:"cache"`
	Runtime string `toml:"runtime"`
}

func TomlXDGV0FromXDG(xdg xdg.XDG) *TomlXDGV0 {
	return &TomlXDGV0{
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
