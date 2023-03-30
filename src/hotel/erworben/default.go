package erworben

import "github.com/friedenberg/zit/src/delta/kennung"

func Default() (k *Objekte) {
	k = &Objekte{
		Akte: Akte{
			DefaultTyp: kennung.MustTyp("md"),
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
