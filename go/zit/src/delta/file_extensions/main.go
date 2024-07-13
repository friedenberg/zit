package file_extensions

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

type FileExtensions struct {
	Zettel   string `toml:"zettel"`
	Organize string `toml:"organize"`
	Typ      string `toml:"typ"`
	Etikett  string `toml:"etikett"`
	Kasten   string `toml:"kasten"`
}

func (a FileExtensions) GetFileExtensionForGattung(
	g1 interfaces.GenreGetter,
) string {
	g := genres.Must(g1)

	switch g {
	case genres.Zettel:
		return a.GetFileExtensionZettel()

	case genres.Type:
		return a.GetFileExtensionTyp()

	case genres.Tag:
		return a.GetFileExtensionEtikett()

	case genres.Repo:
		return a.GetFileExtensionKasten()

	default:
		return ""
	}
}

func (a FileExtensions) GetFileExtensionZettel() string {
	return a.Zettel
}

func (a FileExtensions) GetFileExtensionOrganize() string {
	return a.Organize
}

func (a FileExtensions) GetFileExtensionTyp() string {
	return a.Typ
}

func (a FileExtensions) GetFileExtensionEtikett() string {
	return a.Etikett
}

func (a FileExtensions) GetFileExtensionKasten() string {
	return a.Kasten
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
