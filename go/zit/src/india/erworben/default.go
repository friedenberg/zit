package erworben

import (
	"code.linenisgreat.com/zit/src/bravo/erworben_tools"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

func Default(defaultTyp kennung.Typ) (k Akte) {
	k = Akte{
		Defaults: Defaults{
			Typ:       defaultTyp,
			Etiketten: make([]kennung.Etikett, 0),
		},
		FileExtensions: FileExtensions{
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
