package mutable_config_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/options_tools"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

const (
	TypeV0 = builtin_types.ConfigTypeTomlV0
	TypeV1 = builtin_types.ConfigTypeTomlV1
)

type (
	Blob interface {
		interfaces.MutableStoredConfig
		GetDefaults() Defaults
		GetFileExtensions() interfaces.FileExtensionGetter
		GetPrintOptions() options_print.V0
	}

	Defaults interface {
		GetType() ids.Type
		GetTags() quiter.Slice[ids.Tag]
	}
)

func Default(defaultTyp ids.Type) ids.TypedBlob[Blob] {
	return ids.TypedBlob[Blob]{
		Type: builtin_types.DefaultOrPanic(genres.Config),
		Blob: V1{
			Defaults: DefaultsV1{
				Type: defaultTyp,
				Tags: make([]ids.Tag, 0),
			},
			FileExtensions: file_extensions.V1{
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
