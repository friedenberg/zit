package erworben

type FileExtensions struct {
	Zettel   string `toml:"zettel"`
	Organize string `toml:"organize"`
	Typ      string `toml:"typ"`
	Etikett  string `toml:"etikett"`
}

func (a *FileExtensions) Reset(b *FileExtensions) {
	if b == nil {
		a.Zettel = ""
		a.Organize = ""
		a.Typ = ""
		a.Etikett = ""
	} else {
		a.Zettel = b.Zettel
		a.Organize = b.Organize
		a.Typ = b.Typ
		a.Etikett = b.Etikett
	}
}
