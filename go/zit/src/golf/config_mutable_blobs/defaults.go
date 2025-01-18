package config_mutable_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type DefaultsV0 struct {
	Typ       ids.Type  `toml:"typ"`
	Etiketten []ids.Tag `toml:"etiketten"`
}

func (d DefaultsV0) GetType() ids.Type {
	return d.Typ
}

func (d DefaultsV0) GetTags() quiter.Slice[ids.Tag] {
	return quiter.Slice[ids.Tag](d.Etiketten)
}

type DefaultsV1 struct {
	Type ids.Type  `toml:"type"`
	Tags []ids.Tag `toml:"tags"`
}

func (d DefaultsV1) GetType() ids.Type {
	return d.Type
}

func (d DefaultsV1) GetTags() quiter.Slice[ids.Tag] {
	return quiter.Slice[ids.Tag](d.Tags)
}
