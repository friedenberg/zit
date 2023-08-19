package erworben

import "github.com/friedenberg/zit/src/echo/kennung"

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
	}

	return
}
