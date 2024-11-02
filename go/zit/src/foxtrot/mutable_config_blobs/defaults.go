package mutable_config_blobs

import "code.linenisgreat.com/zit/go/zit/src/echo/ids"

type DefaultsV0 struct {
	Typ       ids.Type  `toml:"typ"`
	Etiketten []ids.Tag `toml:"etiketten"`
}

type DefaultsV1 struct {
	Type ids.Type  `toml:"typ"`
	Tags []ids.Tag `toml:"etiketten"`
}
