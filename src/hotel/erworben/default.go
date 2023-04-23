package erworben

import "github.com/friedenberg/zit/src/delta/kennung"

func Default(defaultTyp kennung.Typ) (k Objekte) {
	k = Objekte{
		Akte: Akte{
			DefaultTyp: defaultTyp,
			FileExtensions: FileExtensions{
				Typ:      "typ",
				Zettel:   "zettel",
				Organize: "md",
				Etikett:  "etikett",
				Kasten:   "kasten",
			},
		},
	}

	return
}
