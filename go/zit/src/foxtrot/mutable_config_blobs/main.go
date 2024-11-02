package mutable_config_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/options_tools"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

const (
	TypeV0     = builtin_types.ConfigTypeV0
	TypeV1     = builtin_types.ConfigTypeV1
	TypeLatest = builtin_types.ConfigTypeLatest
)

type Blob interface {
	interfaces.MutableStoredConfig
}

// TODO version
func Default(defaultTyp ids.Type) ids.TypedBlob[Blob] {
	return ids.TypedBlob[Blob]{
		Type: ids.MustType(TypeLatest),
		Blob: V1{
			Defaults: DefaultsV1{
				Type: defaultTyp,
				Tags: make([]ids.Tag, 0),
			},
			FileExtensions: file_extensions.FileExtensions{
				Type:     "type",
				Zettel:   "zettel",
				Organize: "md",
				Tag:      "tag",
				Repo:     "repo",
			},
			Tools: options_tools.Options{
				Merge: []string{
					"vimdiff",
				},
			},
		},
	}
}
