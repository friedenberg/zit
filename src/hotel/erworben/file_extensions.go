package erworben

type FileExtensions struct {
	Zettel   string `toml:"zettel"`
	Organize string `toml:"organize"`
	Typ      string `toml:"typ"`
	Etikett  string `toml:"etikett"`
	Kasten   string `toml:"kasten"`
}

func (a *FileExtensions) Reset() {
	a.Zettel = ""
	a.Organize = ""
	a.Typ = ""
	a.Etikett = ""
	a.Kasten = ""
}

func (a *FileExtensions) ResetWith(b FileExtensions) {
	a.Zettel = b.Zettel
	a.Organize = b.Organize
	a.Typ = b.Typ
	a.Etikett = b.Etikett
	a.Kasten = b.Kasten
}
