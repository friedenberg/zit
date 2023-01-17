package konfig

func Default() (k *Objekte) {
	k = &Objekte{
		Akte: Akte{
			FileExtensions: FileExtensions{
				Typ:      "typ",
				Zettel:   "zettel",
				Organize: "md",
				Etikett:  "etikett",
			},
		},
	}

	return
}
