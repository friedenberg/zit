package mutable_config_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/bravo/options_tools"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

const (
	TypeV0     = "toml-config-v0"
	TypeV1     = "toml-config-v1"
	TypeLatest = TypeV1
)

type Blob interface{}

// TODO version
func Default(defaultTyp ids.Type) (k V0) {
	k = V0{
		Defaults: DefaultsV0{
			Typ:       defaultTyp,
			Etiketten: make([]ids.Tag, 0),
		},
		FileExtensions: file_extensions.FileExtensions{
			Typ:      "typ",
			Zettel:   "zettel",
			Organize: "md",
			Etikett:  "etikett",
			Kasten:   "kasten",
		},
		Tools: options_tools.Options{
			Merge: []string{
				"vimdiff",
			},
		},
	}

	return
}
