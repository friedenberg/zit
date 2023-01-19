package erworben

type FileExtensions struct {
	Zettel   string `toml:"zettel"`
	Organize string `toml:"organize"`
	Typ      string `toml:"typ"`
	Etikett  string `toml:"etikett"`
}

func (a *FileExtensions) Reset() {
	a.Zettel = ""
	a.Organize = ""
	a.Typ = ""
	a.Etikett = ""
}

func (a *FileExtensions) ResetWith(b FileExtensions) {
	a.Zettel = b.Zettel
	a.Organize = b.Organize
	a.Typ = b.Typ
	a.Etikett = b.Etikett
}
