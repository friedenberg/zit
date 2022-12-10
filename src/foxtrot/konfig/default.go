package konfig

func Default() (k *Objekte) {
	k = &Objekte{
		Akte: Toml{
			FileExtensions: FileExtensions{
				Typ:      "typ",
				Zettel:   "zettel",
				Organize: "md",
			},
		},
	}

	return
}
