package file_extensions

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
)

// TODO make new non-german version
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
		return a.GetFileExtensionType()

	case genres.Tag:
		return a.GetFileExtensionTag()

	case genres.Repo:
		return a.GetFileExtensionRepo()

	default:
		return ""
	}
}

func (a FileExtensions) GetFileExtensionGetter() interfaces.FileExtensionGetter {
	return a
}

func (a FileExtensions) GetFileExtensionZettel() string {
	return a.Zettel
}

func (a FileExtensions) GetFileExtensionOrganize() string {
	return a.Organize
}

func (a FileExtensions) GetFileExtensionType() string {
	return a.Typ
}

func (a FileExtensions) GetFileExtensionTag() string {
	return a.Etikett
}

func (a FileExtensions) GetFileExtensionRepo() string {
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
