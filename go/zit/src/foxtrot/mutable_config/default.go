package mutable_config

import (
	"code.linenisgreat.com/zit/go/zit/src/bravo/erworben_tools"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

// TODO version
func Default(defaultTyp ids.Type) (k Blob) {
	k = Blob{
		Defaults: Defaults{
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
		Tools: erworben_tools.Tools{
			Merge: []string{
				"vimdiff",
			},
		},
	}

	return
}
